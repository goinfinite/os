package apiMiddleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func SetDefaultHeaders(apiBasePath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			c.Response().Header().Set(
				"Cache-Control", "no-store, no-cache, must-revalidate",
			)
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set(
				"Access-Control-Allow-Headers",
				"X-Requested-With, Content-Type, Accept, Origin, Authorization",
			)
			c.Response().Header().Set(
				"Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS, DELETE, PUT",
			)

			if req.Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			urlPath := req.URL.Path
			isNotApi := !strings.HasPrefix(urlPath, apiBasePath)

			if isNotApi {
				return next(c)
			}

			if req.Header.Get("Content-Type") == "" {
				req.Header.Set("Content-Type", "application/json")
			}

			c.Response().Header().Set("Content-Type", "application/json")

			return next(c)
		}
	}
}
