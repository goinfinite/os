package apiMiddleware

import (
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
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
			isReadOnlyModeEnabled, err := tkVoUtil.InterfaceToBool(os.Getenv("READ_ONLY_MODE"))
			if err != nil || !isReadOnlyModeEnabled {
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
