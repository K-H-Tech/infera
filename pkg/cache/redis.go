package cache

import (
	"context"
	"errors"
	"zarinpal-platform/core/logger"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	config *RedisConfig
	client *redis.Client
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisCache(config *RedisConfig) Cache {
	rc := &redisCache{
		config: config,
	}
	rc.connect()
	return rc
}

func (rc *redisCache) connect() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     rc.config.Addr,
		Password: rc.config.Password,
		DB:       rc.config.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Log.Fatalf("redis connect err: %v", err)
	}
	rc.client = rdb
	logger.Log.Info("successfully connected to redis")
}

func (rc *redisCache) Set(ctx context.Context, key string, value interface{}, opts ...SetOption) error {
	options := &setOptions{}
	for _, opt := range opts {
		opt(options)
	}

	ttl := options.ttl
	return rc.client.Set(ctx, key, value, ttl).Err()
}

func (rc *redisCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return result, true, nil
}
