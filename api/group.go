package api

import (
	"net/http"
	"strconv"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/labstack/echo"
)

func getGroups(c echo.Context) error {
	backend := c.Get(middleware.BackendKey).(model.Backends)

	var r map[string]model.UserGroup
	var err error
	if len(c.QueryParams()) > 0 {
		for k, v := range c.QueryParams() {
			switch k {
			case "id":
				id, err := strconv.Atoi(v[0])
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, nil)
				}

				r, err = backend.FindGroupByID(id)
				if err != nil {
					return errorResponse(c, err)
				}

			case "name":
				r, err = backend.FindGroupByName(v[0])
				if err != nil {
					return errorResponse(c, err)
				}
			default:
				return echo.NewHTTPError(http.StatusBadRequest, nil)
			}
		}
	} else {
		r, err = backend.Groups()
		if err != nil {
			return errorResponse(c, err)
		}
	}

	return c.JSON(http.StatusOK, toSlice(r))
}

func GroupEndpoints(g *echo.Group) {
	g.GET("/groups", getGroups)
}
