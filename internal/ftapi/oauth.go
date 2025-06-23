package ftapi

import (
	"context"
	"encoding/json"
	"fmt"
	"ftbadge/internal/cache"
	"ftbadge/internal/utils"
	"net/http"
	"strings"
	"time"
)

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	accessTokenCacheKey = "oauth:access-token"
	oauthTokenEndpoint  = "https://api.intra.42.fr/oauth/token"
)

func FetchOAuthToken(ctx context.Context) (*oauthTokenResponse, error) {
	clientID := utils.MustGetEnv("CLIENT_ID")
	clientSecret := utils.MustGetEnv("CLIENT_SECRET")

	data := fmt.Sprintf(
		"grant_type=client_credentials&client_id=%s&client_secret=%s",
		clientID, clientSecret,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthTokenEndpoint, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("could not create OAuth token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending OAuth token request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from OAuth token endpoint: %s", resp.Status)
	}

	var tokenResp oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode OAuth token response: %w", err)
	}

	return &tokenResp, nil
}

func CacheAccessToken(ctx context.Context, cs cache.CacheStore) (string, error) {
	tokenResp, err := FetchOAuthToken(ctx)
	if err != nil {
		return "", fmt.Errorf("could not obtain new OAuth token: %w", err)
	}

	ttl := time.Duration(tokenResp.ExpiresIn) * time.Second
	if err := cs.Set(ctx, accessTokenCacheKey, tokenResp.AccessToken, ttl); err != nil {
		return "", fmt.Errorf("failed to cache access token: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func GetOrCacheAccessToken(ctx context.Context, cs cache.CacheStore) (string, error) {
	cachedValue, found, err := cs.Get(ctx, accessTokenCacheKey)
	if err != nil {
		return "", fmt.Errorf("error retrieving OAuth token from cache: %w", err)
	}
	if !found {
		return CacheAccessToken(ctx, cs)
	}
	return cachedValue, nil
}
