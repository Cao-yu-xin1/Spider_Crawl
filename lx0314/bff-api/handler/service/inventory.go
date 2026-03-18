package service

import (
	"encoding/json"
	"fmt"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/mq1"
	"github.com/Cao-yu-xin1/Spider_Crawl/lx0314/utils"
	"log"
	"sync"
	"time"
)

// InventoryService 库存服务
type InventoryService struct {
	mq        *mq1.QueueProducer
	redis     *utils.RedisClient
	inventory map[string]int
	mu        sync.RWMutex
}

// InventoryMessage 库存消息
type InventoryMessage struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	OrderID   string `json:"order_id"`
}

// NewInventoryService 创建库存服务
func NewInventoryService(mq *mq1.QueueProducer, redis *utils.RedisClient) *InventoryService {
	return &InventoryService{
		mq:        mq,
		redis:     redis,
		inventory: make(map[string]int),
		mu:        sync.RWMutex{},
	}
}

// InitInventory 初始化库存
func (s *InventoryService) InitInventory(productID string, quantity int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inventory[productID] = quantity
	log.Printf("初始化库存: product=%s, quantity=%d", productID, quantity)
}

// DeductInventory 扣减库存（生产者）
func (s *InventoryService) DeductInventory(productID string, quantity int, orderID string) error {
	s.mu.RLock()
	currentStock, exists := s.inventory[productID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("商品不存在: %s", productID)
	}

	if currentStock < quantity {
		return fmt.Errorf("库存不足: 当前库存=%d, 需要=%d", currentStock, quantity)
	}

	// 创建库存消息
	msg := InventoryMessage{
		ProductID: productID,
		Quantity:  quantity,
		OrderID:   orderID,
	}

	// 发送消息到队列
	data, _ := json.Marshal(msg)
	err := s.mq.SendMsg("inventory_deduct", string(data))
	if err != nil {
		return fmt.Errorf("发送库存扣减消息失败: %v", err)
	}

	log.Printf("库存扣减消息已发送: product=%s, quantity=%d, order=%s", productID, quantity, orderID)
	return nil
}

// HandleInventoryDeduct 处理库存扣减（消费者）
func (s *InventoryService) HandleInventoryDeduct() error {
	handler := func(msgBody string) {
		// 解析库存消息
		var inventoryMsg InventoryMessage
		if err := json.Unmarshal([]byte(msgBody), &inventoryMsg); err != nil {
			log.Printf("解析消息失败: %v", err)
			return
		}

		// 幂等性检查：使用Redis SetNX
		lockKey := fmt.Sprintf("inventory:deduct:%s:%s", inventoryMsg.ProductID, inventoryMsg.OrderID)
		success, err := s.redis.SetNX(lockKey, "1", 24*time.Hour)
		if err != nil {
			log.Printf("Redis操作失败: %v", err)
			return
		}

		if !success {
			log.Printf("消息重复消费: product=%s, order=%s", inventoryMsg.ProductID, inventoryMsg.OrderID)
			return
		}

		// 执行库存扣减
		s.mu.Lock()
		currentStock := s.inventory[inventoryMsg.ProductID]
		if currentStock >= inventoryMsg.Quantity {
			s.inventory[inventoryMsg.ProductID] = currentStock - inventoryMsg.Quantity
			log.Printf("库存扣减成功: product=%s, quantity=%d, 剩余库存=%d, order=%s",
				inventoryMsg.ProductID, inventoryMsg.Quantity, s.inventory[inventoryMsg.ProductID], inventoryMsg.OrderID)
		} else {
			log.Printf("库存不足: product=%s, 当前库存=%d, 需要=%d, order=%s",
				inventoryMsg.ProductID, currentStock, inventoryMsg.Quantity, inventoryMsg.OrderID)
		}
		s.mu.Unlock()
	}

	// 订阅库存扣减消息
	return s.mq.SubscribeMsg("inventory_deduct", handler)
}

// GetInventory 获取库存
func (s *InventoryService) GetInventory(productID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.inventory[productID]
}
