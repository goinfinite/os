package main

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//go:embed docs/swagger.json
var swaggerJson []byte

func swaggerRoute(baseRoute *echo.Group) {
	swaggerGroup := baseRoute.Group("/swagger")

	swaggerGroup.GET("/swagger.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, echo.MIMEApplicationJSON, swaggerJson)
	})

	swaggerUrl := echoSwagger.URL("swagger.json")
	swaggerGroup.GET("/*", echoSwagger.EchoWrapHandler(swaggerUrl))
}

func authRoutes(baseRoute *echo.Group) {
}

func RouterInit(baseRoute *echo.Group) {
	swaggerRoute(baseRoute)
	authRoutes(baseRoute)
}
