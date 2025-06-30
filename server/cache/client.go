package cache

import (
	"context"
	"time"
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
