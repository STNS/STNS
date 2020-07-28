package middleware

import (
	"strconv"
	"strings"

	"github.com/STNS/STNS/v2/model"
	"github.com/labstack/echo"
)

func AddHeader(backend model.Backend) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			if strings.Index(c.Path(), "/users") > 0 {
				if backend.HighestUserID() != 0 && backend.LowestUserID() != 0 {
					c.Response().Header().Add("USER-HIGHEST-ID", strconv.Itoa(backend.HighestUserID()))
					c.Response().Header().Add("USER-LOWEST-ID", strconv.Itoa(backend.LowestUserID()))
				}
			} else if strings.Index(c.Path(), "/groups") > 0 {
				if backend.HighestGroupID() != 0 && backend.LowestGroupID() != 0 {
					c.Response().Header().Add("GROUP-HIGHEST-ID", strconv.Itoa(backend.HighestGroupID()))
					c.Response().Header().Add("GROUP-LOWEST-ID", strconv.Itoa(backend.LowestGroupID()))
				}
			}
			if err := next(c); err != nil {
				return err
			}

			return nil
		})
	}
}
