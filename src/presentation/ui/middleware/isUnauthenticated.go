package uiMiddleware

import (
	"net/http"
	"regexp"
	"strings"
)

var unauthenticatedUiCallsRegex *regexp.Regexp = regexp.MustCompile(
	`^/(login|assets|setup|dev)/?`,
)

func IsUnauthenticatedUiCall(httpRequest *http.Request, uiBasePath string) bool {
	isNotUi := !strings.HasPrefix(httpRequest.URL.Path, uiBasePath)
	if isNotUi {
		return true
	}
	uiCallWithoutBasePath := strings.TrimPrefix(httpRequest.URL.Path, uiBasePath)
	return unauthenticatedUiCallsRegex.MatchString(uiCallWithoutBasePath)
}