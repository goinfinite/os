package restApi

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
	restApiController "github.com/speedianet/sam/src/presentation/api/controller"
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
	authGroup := baseRoute.Group("/auth")
	authGroup.POST("/login/", restApiController.AuthLoginController)
}

func userRoutes(baseRoute *echo.Group) {
	userGroup := baseRoute.Group("/user")
	userGroup.POST("/", restApiController.AddUserController)
	userGroup.PUT("/", restApiController.UpdateUserController)
}

func registerRestApiRoutes(baseRoute *echo.Group) {
	swaggerRoute(baseRoute)
	authRoutes(baseRoute)
	userRoutes(baseRoute)
}
