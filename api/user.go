package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/labstack/echo"
)

func getUsers(c echo.Context) error {
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
				return echo.NewHTTPError(http.StatusBadRequest, nil)
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

type PasswordChangeParams struct {
	CurrentPassword string
	NewPassword     string
}

func updateUserPassword(c echo.Context) error {
	backend := c.Get(middleware.BackendKey).(model.Backends)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	u := PasswordChangeParams{}
	if err := c.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	r, err := backend.FindUserByID(id)
	if err != nil {
		return errorResponse(c, err)
	}

	for _, us := range r {
		user := us.(*model.User)
		if user.Password == u.CurrentPassword {
			user.Password = u.CurrentPassword

			err := backend.UpdateUser(user.ID, user)
			if err != nil {
				return errorResponse(c, err)
			}
			return c.JSON(http.StatusOK, user)
		}
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("user id:%d unmatch password", id))
	}
	return echo.NewHTTPError(http.StatusBadRequest, "user notfound")
}

func UserEndpoints(g *echo.Group) {
	g.GET("/users", getUsers)
	g.PUT("/users/password/:id", updateUserPassword)
}
