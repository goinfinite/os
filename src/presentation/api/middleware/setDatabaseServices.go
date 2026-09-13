package apiMiddleware

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
	"github.com/labstack/echo/v4"
)

func SetDatabaseServices(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *tkInfraDb.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("persistentDbSvc", persistentDbSvc)
			c.Set("transientDbSvc", transientDbSvc)
			c.Set("trailDbSvc", trailDbSvc)
			return next(c)
		}
	}
}
