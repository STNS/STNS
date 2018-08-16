package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	"github.com/labstack/echo"
)

func newContext(path string, queryParams map[string]string, config *stns.Config) (echo.Context, *httptest.ResponseRecorder) {
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
	ctx.Set(middleware.ConfigKey, config)
	ctx.SetPath(path)

	return ctx, rec
}

func testConfig() *stns.Config {
	return &stns.Config{
		Users: &model.Users{
			"user1": &model.User{
				Base: model.Base{
					ID:   1,
					Name: "User1",
				},
			},
			"user2": &model.User{
				Base: model.Base{
					ID:   2,
					Name: "User2",
				},
			},
		},
	}
}
