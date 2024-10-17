package ui

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	uiMiddleware "github.com/goinfinite/os/src/presentation/ui/middleware"
	"github.com/labstack/echo/v4"
)

func UiInit(
	e *echo.Echo,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) {
	basePath := ""
	baseRoute := e.Group(basePath)

	e.Use(uiMiddleware.Authentication(persistentDbSvc))

	router := NewRouter(baseRoute, persistentDbSvc, transientDbSvc, trailDbSvc)
	router.RegisterRoutes()
}
