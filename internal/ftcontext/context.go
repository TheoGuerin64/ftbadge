package ftcontext

import (
	"ftbadge/internal/cache"

	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context

	CacheStore cache.CacheStore
}

func New(ctx echo.Context, cs cache.CacheStore) *Context {
	return &Context{
		Context:    ctx,
		CacheStore: cs,
	}
}

func Convert(ctx echo.Context) *Context {
	if c, ok := ctx.(*Context); ok {
		return c
	}
	panic("unable to convert echo.Context to *ftcontext.Context")
}

func Middleware(cs cache.CacheStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newCtx := New(ctx, cs)
			return next(newCtx)
		}
	}
}
