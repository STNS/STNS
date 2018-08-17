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

	b := model.NewBackendTomlFile(config.Users, config.Groups)
	ctx.Set(middleware.BackendKey, b)

	ctx.SetPath(path)

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
