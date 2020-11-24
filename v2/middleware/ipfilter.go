package middleware

import (
	"net"

	"github.com/jpillora/ipfilter"
	"github.com/labstack/echo"
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
			if c.Path() == "/" || c.Path() == "/status" {
				return next(c)
			}

			f := ipfilter.New(ipfilter.Options{
				AllowedIPs:     config.AllowIPs,
				BlockByDefault: true,
			})
			ip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)

			if config.AllowIPs != nil && !f.Allowed(ip) {
				if config.Logger != nil {
					config.Logger.Debugf("access denied %s", c.Request().RemoteAddr)
				}
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
