package apiMiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetDefaultHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		c.Response().Header().Set(
			"Content-Type", "application/json",
		)
		c.Response().Header().Set(
			"Cache-Control", "no-store, no-cache, must-revalidate",
		)
		c.Response().Header().Set(
			"Access-Control-Allow-Origin", "*",
		)
		c.Response().Header().Set(
			"Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Accept, Origin, Authorization",
		)
		c.Response().Header().Set(
			"Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS, DELETE, PUT",
		)

		if c.Request().Method == "OPTIONS" {
			return c.NoContent(http.StatusOK)
		}

		return next(c)
	}
}
