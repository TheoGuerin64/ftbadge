package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"

	"ftbadge/internal/utils"
)

type CacheEntry struct {
	Key   string
	Value string
	TTL   time.Duration
}

type CacheClient interface {
	Get(ctx context.Context, key string) (string, bool, error)
	BulkSet(ctx context.Context, entries []CacheEntry) error
	BulkGet(ctx context.Context, keys ...string) ([]*string, error)
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	redisURL := utils.MustGetEnv("REDIS_URL")
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %q: %w", redisURL, err)
	}
	options.MaintNotificationsConfig = &maintnotifications.Config{
		Mode: maintnotifications.ModeDisabled,
	}

	client := redis.NewClient(options)
	return &RedisClient{client}, nil
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, bool, error) {
	value, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (rc *RedisClient) BulkSet(ctx context.Context, entries []CacheEntry) error {
	pipeline := rc.client.Pipeline()
	for _, entry := range entries {
		pipeline.Set(ctx, entry.Key, entry.Value, entry.TTL)
	}

	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("failed to set values %q in Redis using pipeline: %w", entries, err)
	}
	return nil
}

func (rc *RedisClient) BulkGet(ctx context.Context, keys ...string) ([]*string, error) {
	rawValues, err := rc.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve keys %q from Redis using MGET: %w", keys, err)
	}

	values, err := utils.MapSlice(rawValues, utils.AnyToStringPointer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert MGET result to []*string for keys %q: %w", keys, err)
	}
	return values, nil
}

type LocalClient struct {
	cache *ristretto.Cache[string, string]
}

func NewLocalClient() (*LocalClient, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e5,
		MaxCost:     100 << 20,
		BufferItems: 64,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create local cache: %w", err)
	}

	return &LocalClient{cache}, nil
}

func (lc *LocalClient) Get(ctx context.Context, key string) (string, bool, error) {
	value, found := lc.cache.Get(key)
	if !found {
		return "", false, nil
	}
	return value, true, nil
}

func (lc *LocalClient) BulkSet(ctx context.Context, entries []CacheEntry) error {
	for _, entry := range entries {
		lc.cache.SetWithTTL(entry.Key, entry.Value, int64(len(entry.Value)), entry.TTL)
	}
	lc.cache.Wait()
	return nil
}

func (lc *LocalClient) BulkGet(ctx context.Context, keys ...string) ([]*string, error) {
	values := make([]*string, len(keys))
	for index, key := range keys {
		value, found := lc.cache.Get(key)
		if !found {
			values[index] = nil
			continue
		}

		values[index] = &value
	}
	return values, nil
}
