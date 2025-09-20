package ftapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ftbadge/internal/cache"
	"ftbadge/internal/utils"
)

type userReponse struct {
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

type User struct {
	Email     string
	Name      string
	Role      string
	AvatarURL string
	Grade     string
	Level     float64
	Cursus    string
}

func createUser(userResp *userReponse) *User {
	grade := "N/A"
	level := 0.0
	cursusName := "N/A"

	if len(userResp.CursusUsers) > 0 {
		cursus := userResp.CursusUsers[len(userResp.CursusUsers)-1]
		grade = cursus.Grade
		level = cursus.Level
		cursusName = cursus.Cursus.Name
	}

	return &User{
		Email:     userResp.Email,
		Name:      userResp.Displayname,
		Role:      userResp.Kind,
		AvatarURL: userResp.Image.Versions.Medium,
		Grade:     grade,
		Level:     level,
		Cursus:    cursusName,
	}
}

func (c *Client) GetUser(ctx context.Context, cm *cache.CacheManager, login string) (*User, error) {
	if cachedValue, isCached := cm.Get(cache.CacheKeyUser); isCached {
		var user User
		if err := json.Unmarshal([]byte(cachedValue), &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached user: %w", err)
		}
		return &user, nil
	}

	accessToken, err := c.GetAccessToken(ctx, cm)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %w", err)
	}

	headers := http.Header{}
	headers.Set("Accept-Encoding", "gzip")
	headers.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.Get(ctx, "/users/"+login, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send user request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status from user endpoint: %d %s", resp.StatusCode, resp.Status)
	}

	var data []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		decompressed, err := utils.DecompressGzip(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress user response: %w", err)
		}
		data = decompressed
	} else {
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read user response: %w", err)
		}
		data = bytes
	}

	var userResp userReponse
	if err := json.Unmarshal([]byte(data), &userResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached user: %w", err)
	}
	user := createUser(&userResp)

	cacheData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user for caching: %w", err)
	}
	cm.Set(cache.CacheKeyUser, string(cacheData))

	return user, nil
}
