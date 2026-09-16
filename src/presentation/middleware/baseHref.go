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

			// The deployment proxy strips its prefix from request URLs, so the app
			// sees root-relative paths. The HTML document still needs the prefix to
			// resolve its assets and links. The proxy sends the prefix in the
			// X-Base-Href header.
			//
			// Do not confuse the base href with the base paths. The base paths route
			// requests. The base href serves the document.
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
