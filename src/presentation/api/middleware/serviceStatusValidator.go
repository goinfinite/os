package apiMiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/presentation/shared"
)

func ServiceStatusValidator(serviceNameStr string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := shared.CheckServices(serviceNameStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
					"status": http.StatusBadRequest,
					"body":   err.Error(),
				})
			}

			return next(c)
		}
	}
}
