package ui

import (
	"github.com/labstack/echo/v4"
)

func UiInit(
	e *echo.Echo,
) {
	basePath := "/_"
	baseRoute := e.Group(basePath)

	router := NewRouter(baseRoute)
	router.RegisterRoutes()
}
