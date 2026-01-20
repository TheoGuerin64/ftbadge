package ftapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"ftbadge/internal/cache"
	"ftbadge/internal/utils"
)

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	grantType = "client_credentials"
)

func (c *Client) GetAccessToken(ctx context.Context, cm *cache.CacheManager) (string, error) {
	if cachedValue, isCached := cm.Get(cache.CacheKeyAccessToken); isCached {
		return cachedValue, nil
	}

	clientID := utils.MustGetEnv("FT_CLIENT_ID")
	clientSecret := utils.MustGetEnv("FT_CLIENT_SECRET")

	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	resp, err := c.PostForm(ctx, "/oauth/token", nil, data)
	if err != nil {
		return "", fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status from token endpoint: %d %s", resp.StatusCode, resp.Status)
	}

	var tokenResp oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response from token endpoint: %w", err)
	}

	accessToken := tokenResp.AccessToken
	ttl := time.Duration(tokenResp.ExpiresIn) * time.Second
	if err := cm.SetWithTTL(cache.CacheKeyAccessToken, accessToken, ttl); err != nil {
		return "", fmt.Errorf("failed to cache access token: %w", err)
	}

	return accessToken, nil
}
