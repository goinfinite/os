package apiMiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddTrailingSlash(apiBasePath string) echo.MiddlewareFunc {
	return middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusTemporaryRedirect,
		Skipper: func(c echo.Context) bool {
			return IsSkippableApiCall(c.Request(), apiBasePath)
		},
	})
}
