package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/STNS/STNS/v2/middleware"
	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func tomlContext(path string, queryParams map[string]string, config *stns.Config) (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	q := make(url.Values)
	for k, v := range queryParams {
		q.Set(k, v)
	}

	req, err := http.NewRequest(echo.GET, "/?"+q.Encode(), nil)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	ctx := e.NewContext(req, rec)

	b, _ := model.NewBackendTomlFile(config.Users, config.Groups)
	ctx.Set(middleware.BackendKey, b)
	ctx.SetPath(path)

	return ctx, rec
}

func dummyContext(t *testing.T, reqType, reqPath string, args interface{}) (echo.Context, *httptest.ResponseRecorder) {
	var rp []byte
	var jsonerr error
	if args != nil {
		switch args := args.(type) {
		case string:
			rp = []byte(args)
		default:
			rp, jsonerr = json.Marshal(args)
			assert.NoError(t, jsonerr)
		}
	}
	req := httptest.NewRequest(reqType, reqPath, bytes.NewReader(rp))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(req, rec)

	b, _ := model.NewBackendDummy()
	ctx.Set(middleware.BackendKey, b)
	return ctx, rec
}

func testConfig() *stns.Config {
	return &stns.Config{
		Users: &model.Users{
			"user1": &model.User{
				Base: model.Base{
					ID:   1,
					Name: "user1",
				},
			},
			"user2": &model.User{
				Base: model.Base{
					ID:   2,
					Name: "user2",
				},
			},
		},
		Groups: &model.Groups{
			"group1": &model.Group{
				Base: model.Base{
					ID:   1,
					Name: "group1",
				},
			},
			"group2": &model.Group{
				Base: model.Base{
					ID:   2,
					Name: "group2",
				},
			},
		},
	}
}
