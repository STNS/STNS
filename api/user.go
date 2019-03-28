package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/labstack/echo"
	"github.com/tredoe/osutil/user/crypt"
	_ "github.com/tredoe/osutil/user/crypt/md5_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha256_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha512_crypt"
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
					return c.JSON(http.StatusBadRequest, err)
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
				return c.JSON(http.StatusBadRequest, err)
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
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func updateUserPassword(c echo.Context) (ret error) {
	backend := c.Get(middleware.BackendKey).(model.Backends)
	name := c.Param("name")

	params := PasswordChangeParams{}
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	r, err := backend.FindUserByName(name)
	if err != nil {
		return errorResponse(c, err)
	}

	user := r[name].(*model.User)

	defer func() {
		err := recover()
		if err != nil {
			ret = c.JSON(http.StatusBadRequest, "can't support password hash")
			return
		}
	}()

	cr := crypt.NewFromHash(user.Password)
	if cr.Verify(user.Password, []byte(params.CurrentPassword)) != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("user name :%s unmatch password", name))
	}

	v, err := cr.Generate([]byte(params.NewPassword), []byte{})
	if err != nil {
		return errorResponse(c, err)
	}

	user.Password = string(v)

	err = backend.UpdateUser(user)
	if err != nil {
		return errorResponse(c, err)
	}
	return c.JSON(http.StatusNoContent, user)
}

func UserEndpoints(g *echo.Group) {
	g.GET("/users", getUsers)
	g.PUT("/users/password/:name", updateUserPassword)
}
