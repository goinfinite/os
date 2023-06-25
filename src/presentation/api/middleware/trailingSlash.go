package apiMiddleware

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TrailingSlash(basePath string) echo.MiddlewareFunc {
	trailingSlashSkipRegex := regexp.MustCompile(
		`^` + basePath + `/swagger`,
	)

	return middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusTemporaryRedirect,
		Skipper: func(c echo.Context) bool {
			return trailingSlashSkipRegex.MatchString(c.Request().URL.Path)
		},
	})
}
