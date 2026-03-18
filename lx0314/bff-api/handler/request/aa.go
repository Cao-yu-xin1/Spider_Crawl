package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// MessageDeduplicator 消息去重器
type MessageDeduplicator struct {
	client     *redis.Client
	keyPrefix  string        // key前缀，如 "msg:dedup:"
	expiration time.Duration // key过期时间（必须大于消息处理周期）
}

// NewDeduplicator 创建去重器
func NewDeduplicator(client *redis.Client, expiration time.Duration) *MessageDeduplicator {
	return &MessageDeduplicator{
		client:     client,
		keyPrefix:  "msg:dedup:",
		expiration: expiration,
	}
}

// IsDuplicate 判断消息是否重复
// 返回: true=重复(已处理过), false=新消息(首次处理)
func (d *MessageDeduplicator) IsDuplicate(ctx context.Context, msgID string) (bool, error) {
	key := d.keyPrefix + msgID

	// 使用 SETNX 尝试设置 key
	// 使用 SET NX EX 原子命令（Redis 2.6.12+ 推荐）
	success, err := d.client.SetNX(ctx, key, "1", d.expiration).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx failed: %w", err)
	}

	// success=true: key不存在，设置成功 -> 新消息
	// success=false: key已存在 -> 重复消息
	return !success, nil
}

// MarkProcessed 标记消息已处理完成（可选：延长过期时间）
func (d *MessageDeduplicator) MarkProcessed(ctx context.Context, msgID string, retention time.Duration) error {
	key := d.keyPrefix + msgID
	return d.client.Expire(ctx, key, retention).Err()
}

// Delete 手动删除去重标记（用于重试场景）
func (d *MessageDeduplicator) Delete(ctx context.Context, msgID string) error {
	key := d.keyPrefix + msgID
	return d.client.Del(ctx, key).Err()
}
