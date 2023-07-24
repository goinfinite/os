package api

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
	apiController "github.com/speedianet/sam/src/presentation/api/controller"
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
	authGroup.POST("/login/", apiController.AuthLoginController)
}

func databaseRoutes(baseRoute *echo.Group) {
	databaseGroup := baseRoute.Group("/database")
	databaseGroup.GET("/:dbType/", apiController.GetDatabasesController)
	databaseGroup.POST("/:dbType/", apiController.AddDatabaseController)
	databaseGroup.DELETE(
		"/:dbType/:dbName/",
		apiController.DeleteDatabaseController,
	)
	databaseGroup.POST(
		"/:dbType/:dbName/user/",
		apiController.AddDatabaseUserController,
	)
	databaseGroup.DELETE(
		"/:dbType/:dbName/user/:dbUser/",
		apiController.DeleteDatabaseUserController,
	)
}

func o11yRoutes(baseRoute *echo.Group) {
	o11yGroup := baseRoute.Group("/o11y")
	o11yGroup.GET("/overview/", apiController.O11yOverviewController)
}

func userRoutes(baseRoute *echo.Group) {
	userGroup := baseRoute.Group("/user")
	userGroup.GET("/", apiController.GetUsersController)
	userGroup.POST("/", apiController.AddUserController)
	userGroup.PUT("/", apiController.UpdateUserController)
}

func servicesRoutes(baseRoute *echo.Group) {
	servicesGroup := baseRoute.Group("/services")
	servicesGroup.GET("/", apiController.GetServicesController)
	servicesGroup.PUT("/", apiController.UpdateServiceController)
}

func registerApiRoutes(baseRoute *echo.Group) {
	swaggerRoute(baseRoute)
	authRoutes(baseRoute)
	databaseRoutes(baseRoute)
	o11yRoutes(baseRoute)
	userRoutes(baseRoute)
	servicesRoutes(baseRoute)
}
