package main

import (
	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam-api/src/presentation/api/helper"
	customMiddleware "github.com/speedianet/sam-api/src/presentation/api/middleware"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

//	@title			SamBackend
//	@version		0.0.1
//	@description	SpeediaOS AppManager Backend API
//	@termsOfService	https://speedia.net/tos/

//	@contact.name	Speedia Engineering
//	@contact.url	https://speedia.net/
//	@contact.email	eng+swagger@speedia.net

//	@license.name	Speedia Web Services, LLC © 2023. All Rights Reserved.
//	@license.url	https://speedia.net/tos/

//	@securityDefinitions.apikey Bearer
//	@in header
//	@name Authorization
//	@description Type "Bearer" followed by a space and JWT token.

// @host		localhost:10000
// @BasePath	/v1
func main() {
	helper.CheckEnvs()

	e := echo.New()

	basePath := "/v1"
	baseRoute := e.Group(basePath)

	e.Pre(customMiddleware.TrailingSlash(basePath))
	e.Use(customMiddleware.PanicHandler)
	e.Use(customMiddleware.SetDefaultHeaders)

	RouterInit(baseRoute)

	e.Start(":10000")
}
