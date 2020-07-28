package api

import (
	"net/http"

	"github.com/STNS/STNS/v2/model"
	"github.com/labstack/echo"
)

func toSlice(ug map[string]model.UserGroup) []model.UserGroup {
	if ug == nil {
		return nil
	}
	r := []model.UserGroup{}
	for _, v := range ug {
		r = append(r, v)
	}
	return r
}

func errorResponse(c echo.Context, err error) error {
	switch err.(type) {
	case model.NotFoundError:
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusInternalServerError, err)
}
