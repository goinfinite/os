package apiMiddleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddTrailingSlash(apiBasePath string) echo.MiddlewareFunc {
	return middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusTemporaryRedirect,
		Skipper: func(echoContext echo.Context) bool {
			httpRequestUrl := echoContext.Request().URL.Path
			if strings.HasSuffix(httpRequestUrl, "/") {
				return true
			}

			if httpRequestUrl == apiBasePath+"/swagger" {
				return false
			}

			return IsSkippableApiCall(echoContext.Request(), apiBasePath)
		},
	})
}
