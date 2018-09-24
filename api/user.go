package api

import (
	"net/http"
	"strconv"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/labstack/echo"
)

func getUsers(c echo.Context) error {
	backend := c.Get(middleware.BackendKey).(model.GetterBackends)

	var r map[string]model.UserGroup
	var err error
	if len(c.QueryParams()) > 0 {
		for k, v := range c.QueryParams() {
			switch k {
			case "id":
				id, err := strconv.Atoi(v[0])
				if err != nil {
					return c.JSON(http.StatusBadRequest, nil)
				}

				r, err = backend.FindUserByID(id)
				if err != nil {
					return errorResponse(c, err)
				}
			case "name":
				r, err = backend.FindUserByName(v[0])
				if err != nil {
					return errorResponse(c, err)
				}
			default:
				return c.JSON(http.StatusBadRequest, nil)
			}
		}
	} else {
		r, err = backend.Users()
		if err != nil {
			return errorResponse(c, err)
		}
	}
	return c.JSON(http.StatusOK, toSlice(r))
}

func UserEndpoints(g *echo.Group) {
	g.GET("/users", getUsers)
}
