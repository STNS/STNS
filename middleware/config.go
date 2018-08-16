package middleware

import (
	"github.com/STNS/STNS/stns"
	"github.com/labstack/echo"
)

const (
	ConfigKey = "Config"
)

func Config(config *stns.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.Set(ConfigKey, config)
			return next(c)
		})
	}
}
