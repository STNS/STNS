package middleware

import (
	"fmt"
	"net"
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

// NewIPExtractor は trustedProxies が空の場合は接続元アドレスのみを使い、
// 指定がある場合はそのプロキシ経由の X-Forwarded-For のみを解釈する。
func NewIPExtractor(trustedProxies []string) (echo.IPExtractor, error) {
	if len(trustedProxies) == 0 {
		return echo.ExtractIPDirect(), nil
	}

	opts := []echo.TrustOption{
		echo.TrustLoopback(false),
		echo.TrustLinkLocal(false),
		echo.TrustPrivateNet(false),
	}
	for _, p := range trustedProxies {
		ipnet, err := parseTrustedProxy(p)
		if err != nil {
			return nil, err
		}
		opts = append(opts, echo.TrustIPRange(ipnet))
	}
	return echo.ExtractIPFromXFFHeader(opts...), nil
}

func parseTrustedProxy(s string) (*net.IPNet, error) {
	if strings.Contains(s, "/") {
		_, ipnet, err := net.ParseCIDR(s)
		if err != nil {
			return nil, fmt.Errorf("invalid trusted proxy %q: %w", s, err)
		}
		return ipnet, nil
	}

	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("invalid trusted proxy %q", s)
	}
	if v4 := ip.To4(); v4 != nil {
		return &net.IPNet{IP: v4, Mask: net.CIDRMask(32, 32)}, nil
	}
	return &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}, nil
}
