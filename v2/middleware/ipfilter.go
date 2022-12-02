package middleware

import (
	"strings"

	"github.com/jpillora/ipfilter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	IPFilterConfig struct {
		AllowIPs []string
		Logger   *log.Logger
	}
)

func IPFilterWithConfig(config IPFilterConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/" || strings.HasSuffix(c.Path(), "/status") {
				return next(c)
			}

			f := ipfilter.New(ipfilter.Options{
				AllowedIPs:     config.AllowIPs,
				BlockByDefault: true,
			})

			if config.AllowIPs != nil && !f.Allowed(c.RealIP()) {
				if config.Logger != nil {
					config.Logger.Infof("access denied %s", c.RealIP())
				}
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
