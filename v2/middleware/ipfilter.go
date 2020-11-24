package middleware

import (
	"github.com/jpillora/ipfilter"
	"github.com/labstack/echo"
)

type (
	IPFilterConfig struct {
		AllowIPs []string
	}
)

func IPFilterWithConfig(config IPFilterConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/" || c.Path() == "/status" {
				return next(c)
			}

			f := ipfilter.New(ipfilter.Options{
				AllowedIPs:     config.AllowIPs,
				BlockByDefault: true,
			})
			if config.AllowIPs != nil && !f.Allowed(c.Request().RemoteAddr) {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
