package middleware

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// TokenAuthConfig defines the config for TokenAuth middleware.
	TokenAuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Validator is a function to validate TokenAuth credentials.
		// Required.
		Validator TokenAuthValidator
	}

	// TokenAuthValidator defines a function to validate TokenAuth credentials.
	TokenAuthValidator func(string) bool
)

const (
	token = "Token"
)

// TokenAuthWithConfig returns an TokenAuth middleware with config.
// See `TokenAuth()`.
func TokenAuthWithConfig(config TokenAuthConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Validator == nil {
		panic("token-auth middleware requires validator function")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			t := c.Request().Header.Get("Authorization")
			if config.Validator(strings.TrimSpace(strings.Replace(t, "token", "", 1))) {
				return next(c)
			}
			// Need to return `401` for browsers to pop-up login box.
			c.Response().Header().Set(echo.HeaderWWWAuthenticate, t+" realm=Restricted")
			return echo.ErrUnauthorized
		}
	}
}
