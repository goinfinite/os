package apiMiddleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AddTrailingSlash(apiBasePath string) echo.MiddlewareFunc {
	return func(subsequentHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			requestPath := echoContext.Request().URL.Path
			if strings.HasSuffix(requestPath, "/") {
				return subsequentHandler(echoContext)
			}

			isSwaggerRoot := requestPath == apiBasePath+"/swagger"
			shouldNormalizePath := isSwaggerRoot ||
				!IsSkippableApiCall(echoContext.Request(), apiBasePath)
			if !shouldNormalizePath {
				return subsequentHandler(echoContext)
			}

			baseHref := "/"
			if contextBaseHref, assertOk := echoContext.Get("baseHref").(string); assertOk {
				baseHref = contextBaseHref
			}

			slashTerminatedPath := baseHref + strings.TrimPrefix(requestPath, "/") + "/"
			queryString := echoContext.QueryString()
			if queryString != "" {
				slashTerminatedPath += "?" + queryString
			}

			echoContext.Response().Header().Set("Cache-Control", cacheControlNoCacheHeaderValue)

			return echoContext.Redirect(http.StatusTemporaryRedirect, slashTerminatedPath)
		}
	}
}
