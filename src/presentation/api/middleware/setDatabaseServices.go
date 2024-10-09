package apiMiddleware

import (
	"github.com/labstack/echo/v4"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

func SetDatabaseServices(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
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
