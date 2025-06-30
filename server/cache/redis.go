package cache

import (
	"context"
	"fmt"
	"ftbadge/server/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var (
	redisURL = utils.MustGetEnv("REDIS_URL")
)

func NewRedisClient() (*RedisClient, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %q: %w", redisURL, err)
	}

	client := redis.NewClient(opt)
	return &RedisClient{client}, nil
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, bool, error) {
	t := time.Now()
	value, err := rc.client.Get(ctx, key).Result()
	fmt.Printf("Get took %s for key %q\n", time.Since(t), key)
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (rc *RedisClient) BulkSet(ctx context.Context, entries []CacheEntry) error {
	t := time.Now()
	pipeline := rc.client.Pipeline()
	for _, entry := range entries {
		pipeline.Set(ctx, entry.Key, entry.Value, entry.TTL)
	}

	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set values %q in Redis using pipeline: %w", entries, err)
	}
	fmt.Printf("BulkSet took %s for %d entries\n", time.Since(t), len(entries))
	return nil
}

func (rc *RedisClient) BulkGet(ctx context.Context, keys ...string) ([]*string, error) {
	t := time.Now()
	rawValues, err := rc.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve keys %q from Redis using MGET: %w", keys, err)
	}

	values, err := utils.MapSlice(rawValues, utils.AnyToStringPointer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MGET result to []*string for keys %q: %w", keys, err)
	}
	fmt.Printf("BulkGet took %s for %d entries\n", time.Since(t), len(keys))
	return values, nil
}
