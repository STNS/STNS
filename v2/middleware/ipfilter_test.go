package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestIPFilterWithConfig(t *testing.T) {
	type args struct {
		config IPFilterConfig
	}
	tests := []struct {
		name       string
		remoteAddr string
		header     map[string]string
		path       string
		args       args
		want       echo.MiddlewareFunc
		code       int
	}{
		{
			name: "ok",
			args: args{
				config: IPFilterConfig{
					AllowIPs: []string{"1.1.1.1"},
				},
			},
			path:       "/",
			remoteAddr: "1.1.1.1:10000",
			code:       http.StatusOK,
		},
		{
			name: "unmatch ip",
			args: args{
				config: IPFilterConfig{
					AllowIPs: []string{"2.2.2.2"},
				},
			},
			path:       "/users",
			remoteAddr: "1.1.1.1:10000",
			code:       http.StatusUnauthorized,
		},
		{
			name: "xff from untrusted peer",
			args: args{
				config: IPFilterConfig{
					AllowIPs: []string{"2.2.2.2"},
				},
			},
			path:       "/users",
			remoteAddr: "1.1.1.1:10000",
			header:     map[string]string{"X-Forwarded-For": "2.2.2.2"},
			code:       http.StatusUnauthorized,
		},
		{
			name: "status ok",
			args: args{
				config: IPFilterConfig{
					AllowIPs: []string{"2.2.2.2"},
				},
			},
			path:       "/status",
			remoteAddr: "1.1.1.1:10000",
			code:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest(echo.GET, tt.path, nil)
		if err != nil {
			t.Error(err)
		}
		req.RemoteAddr = tt.remoteAddr
		for k, v := range tt.header {
			req.Header.Set(k, v)
		}

		e := echo.New()
		extractor, err := NewIPExtractor(nil)
		if err != nil {
			t.Fatal(err)
		}
		e.IPExtractor = extractor
		ctx := e.NewContext(req, rec)
		ctx.SetPath(tt.path)

		dummy := func(c echo.Context) error {
			return nil
		}

		t.Run(tt.name, func(t *testing.T) {
			ret := IPFilterWithConfig(tt.args.config)(dummy)(ctx)

			if tt.code != http.StatusOK {
				if ret.(*echo.HTTPError).Code != tt.code {
					t.Errorf("unmatch status conde want:%d got:%d", tt.code, ret.(*echo.HTTPError).Code)
				}
			} else {
				if ret != nil {
					t.Errorf("unmatch status conde got:%s", ret)
				}
			}
		})
	}
}

func TestNewIPExtractor(t *testing.T) {
	tests := []struct {
		name           string
		trustedProxies []string
		remoteAddr     string
		xff            string
		want           string
		wantErr        bool
	}{
		{
			name:       "no trusted proxies ignores xff",
			remoteAddr: "1.1.1.1:10000",
			xff:        "2.2.2.2",
			want:       "1.1.1.1",
		},
		{
			name:           "xff from trusted proxy",
			trustedProxies: []string{"10.0.0.0/8"},
			remoteAddr:     "10.0.0.1:10000",
			xff:            "3.3.3.3, 2.2.2.2, 10.0.0.2",
			want:           "2.2.2.2",
		},
		{
			name:           "xff from untrusted peer",
			trustedProxies: []string{"10.0.0.0/8"},
			remoteAddr:     "192.168.0.1:10000",
			xff:            "2.2.2.2",
			want:           "192.168.0.1",
		},
		{
			name:           "single ip",
			trustedProxies: []string{"10.0.0.1"},
			remoteAddr:     "10.0.0.1:10000",
			xff:            "2.2.2.2",
			want:           "2.2.2.2",
		},
		{
			name:           "single ip does not trust neighbor",
			trustedProxies: []string{"10.0.0.1"},
			remoteAddr:     "10.0.0.2:10000",
			xff:            "2.2.2.2",
			want:           "10.0.0.2",
		},
		{
			name:           "single ipv6",
			trustedProxies: []string{"2001:db8::1"},
			remoteAddr:     "[2001:db8::1]:10000",
			xff:            "2.2.2.2",
			want:           "2.2.2.2",
		},
		{
			name:           "invalid ip",
			trustedProxies: []string{"invalid"},
			wantErr:        true,
		},
		{
			name:           "invalid cidr",
			trustedProxies: []string{"10.0.0.0/33"},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, err := NewIPExtractor(tt.trustedProxies)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(echo.GET, "/users", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xff != "" {
				req.Header.Set(echo.HeaderXForwardedFor, tt.xff)
			}
			if got := extractor(req); got != tt.want {
				t.Errorf("unmatch ip want:%s got:%s", tt.want, got)
			}
		})
	}
}
