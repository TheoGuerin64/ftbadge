package handlers

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"ftbadge/internal/cache"
	"ftbadge/internal/ftapi"
	"ftbadge/internal/ftcontext"
	"ftbadge/templates"
	"math"
	"net/http"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	profileCacheTTL = time.Hour
)

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
	Login string `param:"login" validate:"required,alpha"`
}

func renderProfile(ctx context.Context, cs cache.CacheStore, login string) ([]byte, error) {
	user, err := ftapi.GetOrCacheUser(ctx, cs, login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %q: %w", login, err)
	}

	if len(user.CursusUsers) == 0 {
		return nil, fmt.Errorf("user %q has no cursus information", login)
	}
	currentCursusUser := user.CursusUsers[len(user.CursusUsers)-1]
	level, experience := math.Modf(currentCursusUser.Level)

	image, err := ftapi.GetOrCacheImageBase64(ctx, cs, user.Image.Versions.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to execute profile template: %w", err)
	}

	profile := Profile{
		Avatar:     image,
		Name:       user.Displayname,
		Email:      user.Email,
		Role:       user.Kind,
		Cursus:     currentCursusUser.Cursus.Name,
		Grade:      currentCursusUser.Grade,
		Level:      level,
		Experience: experience,
	}

	tmpl := template.Must(template.New("profile").Parse(templates.Profile))
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, profile); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

func getOrCacheProfile(ctx context.Context, cs cache.CacheStore, login string) ([]byte, error) {
	cacheKey := "profile:" + login

	cachedValue, found, err := cs.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve profile from cache for login %q: %w", login, err)
	}
	if found {
		return []byte(cachedValue), nil
	}

	data, err := renderProfile(ctx, cs, login)
	if err != nil {
		return nil, fmt.Errorf("failed to render profile for login %q: %w", login, err)
	}

	if err := cs.Set(ctx, cacheKey, string(data), profileCacheTTL); err != nil {
		return nil, fmt.Errorf("failed to cache profile for login %q: %w", login, err)
	}

	return data, nil
}

func setClientCacheHeaders(ctx *ftcontext.Context, data []byte) error {
	hash := md5.Sum([]byte(data))
	etag := fmt.Sprintf("\"%x\"", hash)
	clientETag := ctx.Request().Header.Get("If-None-Match")

	ctx.Response().Header().Add("Cache-Control", "public, max-age=3600, s-maxage=86400")
	ctx.Response().Header().Add("Etag", etag)
	if clientETag == etag {
		return ctx.NoContent(http.StatusNotModified)
	}

	return nil
}

func ProfileHandler(ec echo.Context) error {
	ctx := ftcontext.Convert(ec)

	param := profileParam{}
	if err := ctx.Bind(&param); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameters").SetInternal(err)
	}
	if err := ctx.Validate(param); err != nil {
		return err
	}

	data, err := getOrCacheProfile(ctx.Request().Context(), ctx.CacheStore, param.Login)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render profile").SetInternal(err)
	}
	ctx.Response().Header().Add("Content-Type", "image/svg+xml")

	if err := setClientCacheHeaders(ctx, data); err != nil {
		return err
	}
	return ctx.XMLBlob(http.StatusOK, data)
}
