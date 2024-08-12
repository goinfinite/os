package ui

import (
	"github.com/labstack/echo/v4"
	uiMiddleware "github.com/speedianet/os/src/presentation/ui/middleware"
)

func UiInit(
	e *echo.Echo,
) {
	basePath := ""
	baseRoute := e.Group(basePath)

	e.Use(uiMiddleware.Authentication())

	router := NewRouter(baseRoute)
	router.RegisterRoutes()
}
