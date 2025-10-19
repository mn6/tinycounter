package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	IncrementCounter(ctx context.Context, user string) (int64, error)
	GetCounter(ctx context.Context, user string) (int64, error)
	SetImage(ctx context.Context, user, style string, imgData []byte, ttl time.Duration) error
	GetImage(ctx context.Context, user, style string) ([]byte, error)
	TrySetCooldown(ctx context.Context, user, ip string, ttl time.Duration) (bool, error)
	Close() error
}

type RedisStore struct {
	client *redis.Client
	prefix string
}

// Redis keys
const (
	RedisKeyCounters = "counters"
	RedisKeyImages   = "images"
)

func NewRedisStore(redisAddr string, redisPass string, redisDb int, redisPrefix string) (*RedisStore, error) {
	opt := &redis.Options{
		Addr:        redisAddr,
		Password:    redisPass,
		DB:          redisDb,
		DialTimeout: 5 * time.Second,
	}
	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	if redisPrefix == "" {
		redisPrefix = "tinycounter"
	}
	return &RedisStore{client: client, prefix: redisPrefix}, nil
}

func (r *RedisStore) hIncrBy(ctx context.Context, key string, field string, incr int64) (int64, error) {
	query := fmt.Sprintf("%s:%s", r.prefix, key)
	val, err := r.client.HIncrBy(ctx, query, field, incr).Result()
	if err != nil {
		return 0, fmt.Errorf("redis hincrby: %w", err)
	}
	return val, nil
}

func (r *RedisStore) hGet(ctx context.Context, key string, field string) (int64, error) {
	query := fmt.Sprintf("%s:%s", r.prefix, key)
	val, err := r.client.HGet(ctx, query, field).Int64()
	if err != nil {
		return 0, fmt.Errorf("redis hget: %w", err)
	}
	return val, nil
}

func (r *RedisStore) IncrementCounter(ctx context.Context, user string) (int64, error) {
	return r.hIncrBy(ctx, RedisKeyCounters, user, 1)
}

func (r *RedisStore) GetCounter(ctx context.Context, user string) (int64, error) {
	return r.hGet(ctx, RedisKeyCounters, user)
}

func (r *RedisStore) SetImage(ctx context.Context, user string, style string, imgData []byte, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s:%s:%s", r.prefix, RedisKeyImages, user, style)
	return r.client.Set(ctx, key, imgData, ttl).Err()
}

func (r *RedisStore) GetImage(ctx context.Context, user string, style string) ([]byte, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", r.prefix, RedisKeyImages, user, style)
	return r.client.Get(ctx, key).Bytes()
}

func (r *RedisStore) TrySetCooldown(ctx context.Context, user, ip string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("%s:cooldown:%s:%s", r.prefix, user, ip)
	set, err := r.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx: %w", err)
	}
	return set, nil
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}
