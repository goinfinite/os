package presentationMiddleware

import (
	tkPresentationMiddleware "github.com/goinfinite/tk/src/presentation/middleware"
	"github.com/labstack/echo/v4"
)

func PanicHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return tkPresentationMiddleware.ApiPanicHandler(next)
}
