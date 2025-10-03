package apiMiddleware

import (
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func IsSkippableApiCall(req *http.Request, apiBasePath string) bool {
	urlPath := req.URL.Path
	isNotApi := !strings.HasPrefix(urlPath, apiBasePath)
	if isNotApi {
		return true
	}

	urlSkipRegex := regexp.MustCompile(
		`^` + apiBasePath + `/(v\d{1,2}/(auth|health|setup)|/swagger)/?`,
	)
	return urlSkipRegex.MatchString(urlPath)
}

func ReadOnlyMode(apiBasePath string) echo.MiddlewareFunc {
	return func(subsequentHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			rawReadOnlyModeEnvVar := os.Getenv("READ_ONLY_MODE")
			if rawReadOnlyModeEnvVar == "" {
				return subsequentHandler(echoContext)
			}

			isReadOnlyModeEnabled, err := strconv.ParseBool(rawReadOnlyModeEnvVar)
			if err != nil {
				return subsequentHandler(echoContext)
			}

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
				return echoContext.JSON(http.StatusLocked, map[string]interface{}{
					"status": http.StatusLocked,
					"body":   "ReadOnlyModeEnabled",
				})
			}

			return subsequentHandler(echoContext)
		}
	}
}
