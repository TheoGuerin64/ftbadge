package ftapi

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"ftbadge/internal/cache"
	"net/http"
	"time"
)

const (
	UserVersion  = 1
	UserCacheTTL = 24 * time.Hour
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Displayname string `json:"displayname"`
	Kind        string `json:"kind"`
	Image       struct {
		Versions struct {
			Medium string `json:"medium"`
		}
	}
	CursusUsers []struct {
		Grade  string  `json:"grade"`
		Level  float64 `json:"level"`
		Cursus struct {
			Name string `json:"name"`
		} `json:"cursus"`
	} `json:"cursus_users"`
}

func GetUser(ctx context.Context, cs cache.CacheStore, login string) (*User, error) {
	accessToken, err := getOrCacheAccessToken(ctx, cs)
	if err != nil {
		return nil, fmt.Errorf("could not get access token: %w", err)
	}

	url := fmt.Sprintf("https://api.intra.42.fr/v2/users/%s", login)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create user request: %w", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending user request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from user endpoint: %s", resp.Status)
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader for user %q response: %w", login, err)
	}
	defer gzipReader.Close()

	var user User
	if err := json.NewDecoder(gzipReader).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &user, nil
}

func GetOrCacheUser(ctx context.Context, cs cache.CacheStore, login string) (*User, error) {
	cacheKey := fmt.Sprintf("user:%d:%s", UserVersion, login)

	cachedValue, found, err := cs.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("could not get user from cache: %w", err)
	}

	if !found {
		user, err := GetUser(ctx, cs, login)
		if err != nil {
			return nil, fmt.Errorf("could not fetch user from API: %w", err)
		}

		data, err := json.Marshal(user)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal user for caching: %w", err)
		}

		if err := cs.Set(ctx, cacheKey, string(data), UserCacheTTL); err != nil {
			return nil, fmt.Errorf("failed to cache user: %w", err)
		}

		return user, nil
	}

	var user User
	if err := json.Unmarshal([]byte(cachedValue), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached user: %w", err)
	}

	return &user, nil
}
