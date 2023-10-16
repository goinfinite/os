package apiMiddleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
)

func PreventServiceDown(serviceNameStr string) echo.MiddlewareFunc {
	servicesQueryRepo := infra.ServicesQueryRepo{}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			serviceName, err := valueObject.NewServiceName(serviceNameStr)

			currentSvcStatus, err := servicesQueryRepo.GetByName(serviceName)
			if err != nil {
				return err
			}

			var badRequestMessage string

			isStopped := currentSvcStatus.Status.String() == "stopped"
			if isStopped {
				badRequestMessage = "ServiceStopped"
			}
			isUninstalled := currentSvcStatus.Status.String() == "uninstalled"
			if isUninstalled {
				badRequestMessage = "ServiceNotInstalled"
			}
			shouldInstall := isStopped || isUninstalled
			if shouldInstall {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
					"status": http.StatusBadRequest,
					"body":   badRequestMessage,
				})
			}

			return next(c)
		}
	}
}
