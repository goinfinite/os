package apiMiddleware

import (
	"net/http"
	"regexp"
	"slices"
	"strings"

	infraHelper "github.com/goinfinite/os/src/infra/helper"
	"github.com/labstack/echo/v4"
)

var skippableApiCallsRegex *regexp.Regexp = regexp.MustCompile(
	`^(/v\d{1,2}/(auth|health|setup)/?|/swagger/?)`,
)

func IsSkippableApiCall(httpReq *http.Request, apiBasePath string) bool {
	isNotApi := !strings.HasPrefix(httpReq.URL.Path, apiBasePath)
	if isNotApi {
		return true
	}
	apiCallWithoutBasePath := strings.TrimPrefix(httpReq.URL.Path, apiBasePath)

	return skippableApiCallsRegex.MatchString(apiCallWithoutBasePath)
}

func ReadOnlyMode(apiBasePath string) echo.MiddlewareFunc {
	return func(subsequentHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			isReadOnlyModeEnabled := infraHelper.IsReadOnlyMode()
			if !isReadOnlyModeEnabled {
				return subsequentHandler(echoContext)
			}

			shouldSkip := IsSkippableApiCall(echoContext.Request(), apiBasePath)
			if shouldSkip {
				return subsequentHandler(echoContext)
			}

			reqMethod := echoContext.Request().Method
			allowedMethods := []string{"GET", "HEAD", "OPTIONS"}
			if !slices.Contains(allowedMethods, reqMethod) {
				return echoContext.JSON(http.StatusLocked, map[string]any{
					"status": http.StatusLocked,
					"body":   "ReadOnlyModeEnabled",
				})
			}

			return subsequentHandler(echoContext)
		}
	}
}
