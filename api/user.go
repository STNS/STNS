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
	if len(c.QueryParams()) > 0 {
		for k, v := range c.QueryParams() {
			switch k {
			case "id":
				id, err := strconv.Atoi(v[0])
				if err != nil {
					return c.JSON(http.StatusBadRequest, nil)
				}

				r = backend.FindUserByID(id)
			case "name":
				r = backend.FindUserByName(v[0])
			default:
				return c.JSON(http.StatusBadRequest, nil)
			}
		}
	} else {
		r = backend.Users()
	}

	if len(r) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, toSlice(r))
}

func UserEndpoints(g *echo.Group) {
	g.GET("/users", getUsers)
}
