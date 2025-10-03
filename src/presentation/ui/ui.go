package ui

import (
	"embed"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	uiMiddleware "github.com/goinfinite/os/src/presentation/ui/middleware"
	"github.com/labstack/echo/v4"
)

//go:embed assets/*
var assetsFiles embed.FS

func UiInit(
	echoInstance *echo.Echo,
	uiBasePath string,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) {
	baseRoute := echoInstance.Group(uiBasePath)

	echoInstance.Use(uiMiddleware.Authentication(uiBasePath, persistentDbSvc))
	echoInstance.Use(uiMiddleware.Embed([]uiMiddleware.EmbedKeyFs{
		{EmbedKey: "assets", EmbedFs: assetsFiles},
	}))

	router := NewRouter(baseRoute, persistentDbSvc, transientDbSvc, trailDbSvc)
	router.RegisterRoutes()
}
