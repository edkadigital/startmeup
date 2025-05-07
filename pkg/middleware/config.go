package middleware

import (
	"github.com/edkadigital/startmeup/config"
	"github.com/edkadigital/startmeup/pkg/context"
	"github.com/labstack/echo/v4"
)

// Config stores the configuration in the request so it can be accessed by the ui.
func Config(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(context.ConfigKey, cfg)
			return next(ctx)
		}
	}
}
