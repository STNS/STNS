package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

func TestIPFilterWithConfig(t *testing.T) {
	type args struct {
		config IPFilterConfig
	}
	tests := []struct {
		name       string
		remoteAddr string
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

		e := echo.New()
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
