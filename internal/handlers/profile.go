package handlers

import (
	"bytes"
	"context"
	"fmt"
	"ftbadge/internal/cache"
	"ftbadge/internal/ftapi"
	"ftbadge/internal/ftcontext"
	"ftbadge/templates"
	"math"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

const (
	profileWidth  = 340
	profileHeight = 140
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
	Width      float64
	Height     float64
}

type profileParam struct {
	Login string  `param:"login" validate:"required,alpha"`
	Scale float64 `query:"scale" validate:"number,min=1,max=10"`
}

func renderProfile(ctx context.Context, cs cache.CacheStore, param profileParam) ([]byte, error) {
	user, err := ftapi.GetOrCacheUser(ctx, cs, param.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %q: %w", param.Login, err)
	}

	if len(user.CursusUsers) == 0 {
		return nil, fmt.Errorf("user %q has no cursus information", param.Login)
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
		Width:      profileWidth * param.Scale,
		Height:     profileHeight * param.Scale,
	}

	tmpl := template.Must(template.New("profile").Parse(templates.Profile))
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, profile); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

func ProfileHandler(ec echo.Context) error {
	ctx := ftcontext.Convert(ec)

	param := profileParam{
		Scale: 1,
	}
	if err := ctx.Bind(&param); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameters").SetInternal(err)
	}
	if err := ctx.Validate(param); err != nil {
		return err
	}

	data, err := renderProfile(ctx.Request().Context(), ctx.CacheStore, param)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to render profile").SetInternal(err)
	}

	ctx.Response().Header().Add("Content-Type", "image/svg+xml")
	ctx.Response().Header().Add("Cache-Control", "public, max-age=600, s-maxage=3600")

	return ctx.XMLBlob(http.StatusOK, data)
}
