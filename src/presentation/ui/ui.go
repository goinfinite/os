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
	e *echo.Echo,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) {
	basePath := ""
	baseRoute := e.Group(basePath)

	e.Use(uiMiddleware.Authentication(persistentDbSvc))
	e.Use(uiMiddleware.Embed([]uiMiddleware.EmbedKeyFs{
		{EmbedKey: "assets", EmbedFs: assetsFiles},
	}))

	router := NewRouter(baseRoute, persistentDbSvc, transientDbSvc, trailDbSvc)
	router.RegisterRoutes()
}
