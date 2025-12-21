package database

import (
	"context"
	"fmt"
	"time"

	"dbapp/internal/config"
	"dbapp/pkg/logger"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("连接Redis失败: %w", err)
	}

	logger.Info("Redis连接成功")
	return RedisClient, nil
}

func GetRedis() *redis.Client {
	return RedisClient
}

// Cache操作封装
func Set(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

func Delete(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

func Exists(key string) (bool, error) {
	count, err := RedisClient.Exists(ctx, key).Result()
	return count > 0, err
}

func Increment(key string) (int64, error) {
	return RedisClient.Incr(ctx, key).Result()
}

func GetInt(key string) (int, error) {
	val, err := RedisClient.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

