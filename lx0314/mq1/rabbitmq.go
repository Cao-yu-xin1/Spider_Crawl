package mq1

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// InventoryDeductionMessage 库存扣减消息
type InventoryDeductionMessage struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// QueueProducer 封装RabbitMQ队列生产者操作
type QueueProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]amqp.Queue
	mu      sync.RWMutex
}

// NewQueueProducer 创建队列生产者实例
func NewQueueProducer(url string) (*QueueProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("连接RabbitMQ失败: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("创建channel失败: %v", err)
	}

	return &QueueProducer{
		conn:    conn,
		channel: ch,
		queues:  make(map[string]amqp.Queue),
	}, nil
}

// declareQueue 声明队列（线程安全）
func (q *QueueProducer) declareQueue(topic string) (amqp.Queue, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if queue, exists := q.queues[topic]; exists {
		return queue, nil
	}

	queue, err := q.channel.QueueDeclare(
		topic, // 队列名称
		true,  // 持久化
		false, // 自动删除
		false, // 排他性
		false, // 不等待
		nil,   // 参数
	)
	if err != nil {
		return queue, fmt.Errorf("声明队列失败: %v", err)
	}

	q.queues[topic] = queue
	return queue, nil
}

// SendMsg 发送消息到队列
// topic: 队列名称
// msg: 消息内容（字符串）
func (q *QueueProducer) SendMsg(topic string, msg string) error {
	_, err := q.declareQueue(topic)
	if err != nil {
		return err
	}

	msgID := uuid.New().String()

	err = q.channel.Publish(
		"",    // exchange (Simple模式)
		topic, // routing key = 队列名
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(msg),
			DeliveryMode: amqp.Persistent,
			MessageId:    msgID,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	return nil
}

// SubscribeMsg 订阅并消费消息
// topic: 队列名称
// handler: 消息处理函数
func (q *QueueProducer) SubscribeMsg(topic string, handler func(msg string)) error {
	queue, err := q.declareQueue(topic)
	if err != nil {
		return err
	}

	msgs, err := q.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("消费消息失败: %v", err)
	}

	go func() {
		for d := range msgs {
			handler(string(d.Body))
			err := d.Ack(false)
			if err != nil {
				fmt.Printf("消息确认失败: %v\n", err)
			}
		}
	}()

	return nil
}

// Close 关闭连接
func (q *QueueProducer) Close() error {
	if err := q.channel.Close(); err != nil {
		return err
	}
	if err := q.conn.Close(); err != nil {
		return err
	}
	return nil
}

// RedisIdempotency Redis幂等性检查
type RedisIdempotency struct {
	client *redis.Client
}

// NewRedisIdempotency 创建Redis幂等性检查器
func NewRedisIdempotency(addr, password string, db int) *RedisIdempotency {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisIdempotency{client: client}
}

// CheckAndMark 检查消息是否已处理并标记
// key: 消息唯一标识（推荐使用md5(msg)）
// ttl: 过期时间（默认24小时）
func (r *RedisIdempotency) CheckAndMark(key string, ttl time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("Redis检查失败: %v", err)
	}

	return result, nil
}

// GenerateMessageKey 生成消息唯一标识
func GenerateMessageKey(msg string) string {
	hash := md5.Sum([]byte(msg))
	return "mq:processed:" + hex.EncodeToString(hash[:])
}
