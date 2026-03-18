package mq1

import (
	"testing"
	"time"
)

// TestQueueProducerSendMsg 测试发送消息
func TestQueueProducerSendMsg(t *testing.T) {
	producer, err := NewQueueProducer("amqp://caoyuxin:caoyuxin@115.190.154.22:5672/kuteng")
	if err != nil {
		t.Fatalf("创建QueueProducer失败: %v", err)
	}
	defer producer.Close()

	topic := "test_inventory_deduction"
	msg := `{"product_id": 1001, "quantity": 5}`

	err = producer.SendMsg(topic, msg)
	if err != nil {
		t.Errorf("发送消息失败: %v", err)
	}
}

// TestQueueProducerSubscribeMsg 测试订阅消息
func TestQueueProducerSubscribeMsg(t *testing.T) {
	producer, err := NewQueueProducer("amqp://caoyuxin:caoyuxin@115.190.154.22:5672/kuteng")
	if err != nil {
		t.Fatalf("创建QueueProducer失败: %v", err)
	}
	defer producer.Close()

	topic := "test_subscribe"
	receivedMsg := make(chan string, 1)

	// 启动订阅
	err = producer.SubscribeMsg(topic, func(msg string) {
		receivedMsg <- msg
	})
	if err != nil {
		t.Errorf("订阅消息失败: %v", err)
	}

	// 发送测试消息
	testMsg := "test message"
	err = producer.SendMsg(topic, testMsg)
	if err != nil {
		t.Errorf("发送测试消息失败: %v", err)
	}

	// 等待接收消息
	select {
	case msg := <-receivedMsg:
		if msg != testMsg {
			t.Errorf("期望消息: %s, 实际消息: %s", testMsg, msg)
		}
	case <-time.After(5 * time.Second):
		t.Error("等待消息超时")
	}
}

// TestRedisIdempotency 测试Redis幂等性检查
// 注意: Redis未运行时测试会失败，跳过该测试
func TestRedisIdempotency(t *testing.T) {
	redisClient := NewRedisIdempotency("127.0.0.1:6379", "", 0)

	msg := "test message for idempotency"
	key := GenerateMessageKey(msg)

	// 第一次检查，应返回true（未处理）
	result, err := redisClient.CheckAndMark(key, 24*time.Hour)
	if err != nil {
		t.Skipf("Redis连接失败，跳过测试: %v", err)
	}
	if !result {
		t.Error("第一次检查应返回true（表示未处理）")
	}

	// 第二次检查，应返回false（已处理）
	result, err = redisClient.CheckAndMark(key, 24*time.Hour)
	if err != nil {
		t.Skipf("Redis连接失败，跳过测试: %v", err)
	}
	if result {
		t.Error("第二次检查应返回false（表示已处理）")
	}
}

// TestGenerateMessageKey 测试消息唯一标识生成
func TestGenerateMessageKey(t *testing.T) {
	msg := "test message"
	key := GenerateMessageKey(msg)

	expectedPrefix := "mq:processed:"
	if key[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("消息key前缀错误，期望: %s, 实际: %s", expectedPrefix, key[:len(expectedPrefix)])
	}

	// 相同消息应生成相同key
	key2 := GenerateMessageKey(msg)
	if key != key2 {
		t.Errorf("相同消息应生成相同key: %s != %s", key, key2)
	}

	// 不同消息应生成不同key
	key3 := GenerateMessageKey("different message")
	if key == key3 {
		t.Error("不同消息应生成不同key")
	}
}

// TestInventoryDeductionMessage 测试库存扣减消息结构
func TestInventoryDeductionMessage(t *testing.T) {
	msg := InventoryDeductionMessage{
		ProductID: 1001,
		Quantity:  5,
	}

	if msg.ProductID != 1001 {
		t.Errorf("期望ProductID: 1001, 实际: %d", msg.ProductID)
	}

	if msg.Quantity != 5 {
		t.Errorf("期望Quantity: 5, 实际: %d", msg.Quantity)
	}
}
