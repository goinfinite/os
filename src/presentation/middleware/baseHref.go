package presentationMiddleware

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/valueObject"
	"github.com/labstack/echo/v4"
)

func BaseHref(rootBasePath, apiBasePath, uiBasePath string) echo.MiddlewareFunc {
	return func(subsequentHandler echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			echoContext.Set("rootBasePath", rootBasePath)
			echoContext.Set("apiBasePath", apiBasePath)
			echoContext.Set("uiBasePath", uiBasePath)

			// BaseHref is used to set the base href of the HTML document.
			// It is set to the root base path by default, unless X-Base-Href header is set.
			// Base href is only used for assets and API calls from the HTML document.
			// NOTE: Do not confuse this with the base paths. Base paths are used for routing.
			// If the application is behind a reverse proxy, it will receive the request URLs
			// as if it was running on the root of the hostname. On the other hand, the base href
			// won't, that's why X-Base-Href header is used.
			baseHrefStr := rootBasePath
			if len(baseHrefStr) == 0 {
				baseHrefStr = "/"
			}
			baseHrefHasTrailingSlash := baseHrefStr[len(baseHrefStr)-1] == '/'
			if !baseHrefHasTrailingSlash {
				baseHrefStr += "/"
			}

			echoContext.Set("baseHref", baseHrefStr)
			rawBaseHref := echoContext.Request().Header.Get("X-Base-Href")
			if rawBaseHref == "" {
				return subsequentHandler(echoContext)
			}

			baseHref, err := valueObject.NewUrlPath(rawBaseHref)
			if err != nil {
				slog.Debug("InvalidBaseHref", slog.Any("rawBaseHref", rawBaseHref))
				return subsequentHandler(echoContext)
			}

			baseHrefStr = baseHref.String()
			if len(baseHrefStr) == 0 {
				baseHrefStr = "/"
			}
			baseHrefHasTrailingSlash = baseHrefStr[len(baseHrefStr)-1] == '/'
			if !baseHrefHasTrailingSlash {
				baseHrefStr += "/"
			}

			echoContext.Set("baseHref", baseHrefStr)
			return subsequentHandler(echoContext)
		}
	}
}
