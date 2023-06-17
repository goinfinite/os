package restApi

import (
	"github.com/labstack/echo/v4"
	restApiHelper "github.com/speedianet/sam/src/presentation/api/helper"
	restApiMiddleware "github.com/speedianet/sam/src/presentation/api/middleware"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

//	@title			SamApi
//	@version		0.0.1
//	@description	SpeediaOS AppManager API
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
func StartRestApi() {
	restApiHelper.CheckEnvs()

	e := echo.New()

	basePath := "/v1"
	baseRoute := e.Group(basePath)

	e.Pre(restApiMiddleware.TrailingSlash(basePath))
	e.Use(restApiMiddleware.PanicHandler)
	e.Use(restApiMiddleware.SetDefaultHeaders)

	RestApiRouterInit(baseRoute)

	e.Start(":10000")
}
