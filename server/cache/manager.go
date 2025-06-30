package cache

import (
	"context"
	"fmt"
	"time"
)

type CacheManager struct {
	id      string
	client  CacheClient
	data    map[string]string
	pending []CacheEntry
}

const (
	CacheKeyAccessToken = "access-token"
	CacheKeyProfile     = "profile"
	CacheKeyUser        = "user"
	CacheKeyAvatar      = "avatar"
)

var CacheKeys = []string{
	CacheKeyAccessToken,
	CacheKeyProfile,
	CacheKeyUser,
	CacheKeyAvatar,
}

type CacheGroup int

const (
	CacheKeyGroupProfile CacheGroup = iota
	CacheKeyGroupData
)

var preFetchGroups = map[CacheGroup][]string{
	CacheKeyGroupProfile: {CacheKeyProfile},
	CacheKeyGroupData:    {CacheKeyAccessToken, CacheKeyUser, CacheKeyAvatar},
}

func generateAccessTokenKey(id string) string { return "access-token" }
func generateProfileKey(id string) string     { return "profile:" + id }
func generateUserKey(id string) string        { return "user:" + id }
func generateAvatarKey(id string) string      { return "avatar:" + id }

var cacheKeyGenerators = map[string]func(id string) string{
	CacheKeyAccessToken: generateAccessTokenKey,
	CacheKeyProfile:     generateProfileKey,
	CacheKeyUser:        generateUserKey,
	CacheKeyAvatar:      generateAvatarKey,
}

var cacheKeyTTL = map[string]time.Duration{
	CacheKeyProfile: 1 * time.Hour,
	CacheKeyUser:    24 * time.Hour,
	CacheKeyAvatar:  7 * 24 * time.Hour,
}

func NewCacheManager(ctx context.Context, client CacheClient, id string) (*CacheManager, error) {
	data := make(map[string]string, len(CacheKeys))
	var pending []CacheEntry = nil

	return &CacheManager{id, client, data, pending}, nil
}

func (cm *CacheManager) PreFetch(ctx context.Context, group CacheGroup) error {
	keys, exists := preFetchGroups[group]
	if !exists {
		panic(fmt.Sprintf("pre-fetch group %q does not exist", group))
	}

	cacheKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		generator, exists := cacheKeyGenerators[key]
		if !exists {
			panic(fmt.Sprintf("cache key '%s' does not have a corresponding generator function", key))
		}
		cacheKeys = append(cacheKeys, generator(cm.id))
	}

	if len(cacheKeys) > 1 {
		cacheValues, err := cm.client.BulkGet(ctx, cacheKeys...)
		if err != nil {
			return fmt.Errorf("failed to bulk get cache values for pre-fetch group %q: %w", group, err)
		}

		for index, key := range keys {
			if value := cacheValues[index]; value != nil {
				cm.data[key] = *value
			}
		}
	} else {
		value, exists, err := cm.client.Get(ctx, cacheKeys[0])
		if err != nil {
			return fmt.Errorf("failed to get cache value for key %q in pre-fetch group %q: %w", keys[0], group, err)
		}
		if exists {
			cm.data[keys[0]] = value
		}
	}

	return nil
}

func (cm *CacheManager) Get(cacheKey string) (string, bool) {
	value, exists := cm.data[cacheKey]
	return value, exists
}

func (cm *CacheManager) Set(cacheKey string, value string) {
	ttl, exists := cacheKeyTTL[cacheKey]
	if !exists {
		panic(fmt.Sprintf("cache key '%s' does not have a default TTL", cacheKey))
	}

	cm.SetWithTTL(cacheKey, value, ttl)
}

func (cm *CacheManager) SetWithTTL(cacheKey string, value string, ttl time.Duration) {
	if cm.pending == nil {
		cm.pending = make([]CacheEntry, len(CacheKeys))
	}

	generator, exists := cacheKeyGenerators[cacheKey]
	if !exists {
		panic(fmt.Sprintf("cache key '%s' does not have a corresponding generator function", cacheKey))
	}
	key := generator(cm.id)

	entry := CacheEntry{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	cm.pending = append(cm.pending, entry)
}

func (cm *CacheManager) Flush(ctx context.Context) error {
	if cm.pending == nil {
		return nil
	}

	if err := cm.client.BulkSet(ctx, cm.pending); err != nil {
		return fmt.Errorf("failed to bulk set cache entries: %w", err)
	}

	cm.pending = nil
	return nil
}
