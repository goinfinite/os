package apiMiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

func ServiceStatusValidator(servicesNames []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := sharedHelper.CheckServices(servicesNames)
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
