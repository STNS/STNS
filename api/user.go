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

func updateUserPassword(c echo.Context) error {
	backend := c.Get(middleware.BackendKey).(model.Backend)
	u := struct {
		CurrentPassword string
		NewPassword     string
	}{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	r, err := backend.FindUserByID(id)
	if err != nil {
		return errorResponse(c, err)
	}

	for _, us := range r {
		user := us.(*model.User)
		if user.Password == u.CurrentPassword {
			user.Password = u.CurrentPassword

			err := backend.Update(fmt.Sprintf("/users/name/%s", user.Name), user)
			if err != nil {
				return errorResponse(c, err)
			}
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func UserEndpoints(g *echo.Group) {
	g.GET("/users", getUsers)
	g.PUT("/users/password/:id", updateUserPassword)
}
