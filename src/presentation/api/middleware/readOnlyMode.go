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
		`^` + apiBasePath + `(/v\d{1,2}/(auth|health)|/swagger)`,
	)
	return urlSkipRegex.MatchString(urlPath)
}

func ReadOnlyMode(apiBasePath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rawReadOnlyModeEnvVar := os.Getenv("READ_ONLY_MODE")
			if rawReadOnlyModeEnvVar == "" {
				return next(c)
			}

			isReadOnlyModeEnabled, err := strconv.ParseBool(rawReadOnlyModeEnvVar)
			if err != nil {
				return next(c)
			}

			if !isReadOnlyModeEnabled {
				return next(c)
			}

			shouldSkip := IsSkippableApiCall(c.Request(), apiBasePath)
			if shouldSkip {
				return next(c)
			}

			reqMethod := c.Request().Method
			allowedMethods := []string{"GET", "HEAD", "OPTIONS"}
			if !slices.Contains(allowedMethods, reqMethod) {
				return c.JSON(423, map[string]interface{}{
					"status": 423,
					"body":   "ReadOnlyModeEnabled",
				})
			}

			return next(c)
		}
	}
}
