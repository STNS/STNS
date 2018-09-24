package api

import (
	"net/http"
	"strconv"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/labstack/echo"
)

func getGroups(c echo.Context) error {
	backend := c.Get(middleware.BackendKey).(model.GetterBackend)

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

				r, err = backend.FindGroupByID(id)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, err)
				}

			case "name":
				r, err = backend.FindGroupByName(v[0])
				if err != nil {
					return c.JSON(http.StatusInternalServerError, err)
				}
			default:
				return c.JSON(http.StatusBadRequest, nil)
			}
		}
	} else {
		r, err = backend.Groups()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	if len(r) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, toSlice(r))
}

func GroupEndpoints(g *echo.Group) {
	g.GET("/groups", getGroups)
}
