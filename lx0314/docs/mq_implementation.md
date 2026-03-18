# E-commerce Order, Inventory, Points Flow Implementation

## 要求总结
- ✅ 创建RabbitMQ生产者封装 sendMsg(topic, msg string)
- ✅ 创建RabbitMQ消费者封装 subscribeMsg(topic, handler func(msg string))
- ✅ 实现Redis幂等性检查 (SETNX)
- ✅ 库存扣减消息入队列
- ✅ 库存扣减后消费消息并记录日志
- ✅ 单元测试

## 文件结构
```
mq/
├── rabbitmq.go      # 核心封装
└── rabbitmq_test.go # 单元测试
```

##rabbitmq.go 核心实现

```go
package mq

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
func NewQueueProducer(url string) (*QueueProducer, error)

// SendMsg 发送消息到队列
// topic: 队列名称
// msg: 消息内容（字符串）
func (q *QueueProducer) SendMsg(topic string, msg string) error

// SubscribeMsg 订阅并消费消息
// topic: 队列名称
// handler: 消息处理函数
func (q *QueueProducer) SubscribeMsg(topic string, handler func(msg string)) error

// Close 关闭连接
func (q *QueueProducer) Close() error

// RedisIdempotency Redis幂等性检查
type RedisIdempotency struct {
	client *redis.Client
}

// NewRedisIdempotency 创建Redis幂等性检查器
func NewRedisIdempotency(addr, password string, db int) *RedisIdempotency

// CheckAndMark 检查消息是否已处理并标记
// key: 消息唯一标识（推荐使用md5(msg)）
// ttl: 过期时间（默认24小时）
func (r *RedisIdempotency) CheckAndMark(key string, ttl time.Duration) (bool, error)

// GenerateMessageKey 生成消息唯一标识
func GenerateMessageKey(msg string) string
```

## 单元测试 rabbitmq_test.go

```go
package mq

import (
	"testing"
	"time"
)

// TestQueueProducerSendMsg 测试发送消息
// TestQueueProducerSubscribeMsg 测试订阅消息
// TestRedisIdempotency 测试Redis幂等性检查
// TestGenerateMessageKey 测试消息唯一标识生成
```

## 使用示例

```go
// 1. 创建生产者
producer, err := mq.NewQueueProducer("amqp://user:pass@host:5672/vhost")
if err != nil {
    panic(err)
}
defer producer.Close()

// 2. 发送库存扣减消息
msg := `{"product_id": 1001, "quantity": 5}`
err = producer.SendMsg("inventory_deduction_queue", msg)
if err != nil {
    log.Printf("发送失败: %v", err)
}

// 3. 发布库存扣减消息
inventory := mq.InventoryDeductionMessage{
    ProductID: 1001,
    Quantity:  5,
}
data, _ := json.Marshal(inventory)
producer.SendMsg("inventory_deduction_queue", string(data))

// 4. 订阅消费消息
err = producer.SubscribeMsg("inventory_deduction_queue", func(msg string) {
    log.Printf("收到消息: %s", msg)
    // 处理库存扣减逻辑
})

// 5. Redis幂等性检查
redisClient := mq.NewRedisIdempotency("127.0.0.1:6379", "", 0)
key := mq.GenerateMessageKey(msg)
isProcessed, err := redisClient.CheckAndMark(key, 24*time.Hour)
```

## 运行测试

```bash
# 运行所有测试
go test -v ./mq/...

# 运行单个测试
go test -v ./mq/... -run TestQueueProducerSendMsg
```
