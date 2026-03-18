package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisClient Redis客户端
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("连接Redis失败: %v", err)
	}

	return &RedisClient{client: client}, nil
}

// SetNX 设置分布式锁，用于幂等性检查
func (r *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// Get 获取值
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisClient) Del(key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}
