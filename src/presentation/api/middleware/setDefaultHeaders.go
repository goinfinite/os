package apiMiddleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func SetDefaultHeaders(apiBasePath string) echo.MiddlewareFunc {
	return func(subsequentHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			httpRequest := echoContext.Request()

			echoContext.Response().Header().Set(
				"Cache-Control", "no-store, no-cache, must-revalidate",
			)
			echoContext.Response().Header().Set("Access-Control-Allow-Origin", "*")
			echoContext.Response().Header().Set(
				"Access-Control-Allow-Headers",
				"X-Requested-With, Content-Type, Accept, Origin, Authorization",
			)
			echoContext.Response().Header().Set(
				"Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS, DELETE, PUT",
			)

			if httpRequest.Method == "OPTIONS" {
				return echoContext.NoContent(http.StatusOK)
			}

			urlPath := httpRequest.URL.Path
			isNotApi := !strings.HasPrefix(urlPath, apiBasePath)
			if isNotApi {
				return subsequentHandler(echoContext)
			}

			if httpRequest.Header.Get("Content-Type") == "" {
				httpRequest.Header.Set("Content-Type", "application/json")
			}

			echoContext.Response().Header().Set("Content-Type", "application/json")

			return subsequentHandler(echoContext)
		}
	}
}
