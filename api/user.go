package api

import (
	"net/http"
	"strconv"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/labstack/echo"
)

func getUsers(c echo.Context) error {
	backend := c.Get(middleware.BackendKey).(model.Backend)

	var r map[string]model.UserGroup
	for k, v := range c.QueryParams() {
		switch k {
		case "id":
			id, err := strconv.Atoi(v[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, nil)
			}

			r = backend.FindUserByID(id)
		case "name":
			// name
		default:
			return c.JSON(http.StatusBadRequest, nil)
		}
	}

	if len(r) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, toSlice(r))
}

func UserEndpoints(g *echo.Group) {
	g.GET("/users", getUsers)
}
