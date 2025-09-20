package handlers

import (
	"context"
	"crypto/md5" // #nosec G501 -- only used for ETag generation
	"fmt"
	"math"
	"net/http"
	"net/url"
	"text/template"

	"github.com/labstack/echo/v4"

	"ftbadge/internal/cache"
	"ftbadge/internal/ftapi"
	"ftbadge/internal/templates"
	"ftbadge/internal/utils"
)

type UserNotFoundError struct {
	Login string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user %q not found", e.Login)
}

type Profile struct {
	Avatar     string
	Name       string
	Email      string
	Role       string
	Cursus     string
	Grade      string
	Experience float64
	Level      float64
}

type profileParam struct {
	Login string `param:"login" validate:"required,alphanum,max=32"`
}

const (
	apiBaseURL = "https://api.intra.42.fr/v2"
	cdnBaseURL = "https://cdn.intra.42.fr"
)

var (
	profileTemplate = template.Must(template.New("profile").Parse(templates.Profile))
)

func createProfile(user *ftapi.User, avatar string) *Profile {
	level, experience := math.Modf(user.Level)
	experience = max(experience, 0.001) // Ensure experience is never zero to avoid rendering issues

	return &Profile{
		Avatar:     avatar,
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
		Cursus:     user.Cursus,
		Grade:      user.Grade,
		Level:      level,
		Experience: experience,
	}
}

func renderProfile(ctx context.Context, ftc *ftapi.Client, cc cache.CacheClient, login string) ([]byte, error) {
	cm, err := cache.NewCacheManager(ctx, cc, login)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache manager: %w", err)
	}
	if err := cm.PreFetch(ctx, cache.CacheGroupProfile); err != nil {
		return nil, fmt.Errorf("failed to pre-fetch profile cache group: %w", err)
	}
	if cachedProfile, isCached := cm.Get(cache.CacheKeyProfile); isCached {
		return []byte(cachedProfile), nil
	}
	if err := cm.PreFetch(ctx, cache.CacheGroupData); err != nil {
		return nil, fmt.Errorf("failed to pre-fetch data cache group: %w", err)
	}

	user, err := ftc.GetUser(ctx, cm, login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, &UserNotFoundError{Login: login}
	}

	avatarURL, err := url.Parse(user.AvatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user image URL: %w", err)
	}
	avatar, err := ftc.GetAvatar(ctx, cm, avatarURL.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get avatar: %w", err)
	}

	profile := createProfile(user, avatar)
	data, err := utils.RenderTemplate(profileTemplate, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to render profile template: %w", err)
	}

	cm.Set(cache.CacheKeyProfile, string(data))
	if err := cm.Flush(ctx); err != nil {
		return nil, fmt.Errorf("failed to flush cache: %w", err)
	}

	return data, nil
}

func setCacheHeaders(ctx echo.Context, data []byte) error {
	hash := md5.Sum([]byte(data)) // #nosec G401 -- ETag does not need to be cryptographically secure
	etag := fmt.Sprintf("\"%x\"", hash)
	clientETag := ctx.Request().Header.Get("If-None-Match")

	ctx.Response().Header().Add("Cache-Control", "public, s-maxage=3600, stale-while-revalidate=86400")
	ctx.Response().Header().Add("Etag", etag)

	if clientETag == etag {
		return ctx.NoContent(http.StatusNotModified)
	}
	return nil
}

func profileHandler(ctx echo.Context, cc cache.CacheClient) error {
	ctx.Response().Header().Add("Access-Control-Allow-Origin", "*")

	param := profileParam{}
	if err := ctx.Bind(&param); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameters").SetInternal(err)
	}
	if err := ctx.Validate(param); err != nil {
		return err
	}
	if param.Login == "graph" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid login: 'graph' is not allowed")
	}

	ftc := ftapi.NewClient(apiBaseURL, cdnBaseURL)
	data, err := renderProfile(ctx.Request().Context(), ftc, cc, param.Login)
	if err != nil {
		if _, ok := err.(*UserNotFoundError); ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("User %q not found", param.Login))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render profile").SetInternal(err)
	}
	ctx.Response().Header().Add("Content-Type", "image/svg+xml")

	if err := setCacheHeaders(ctx, data); err != nil {
		return err
	}
	return ctx.XMLBlob(http.StatusOK, data)
}

func GetProfileHandler(cc cache.CacheClient) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return profileHandler(ctx, cc)
	}
}
