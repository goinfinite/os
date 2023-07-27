package api

import (
	"github.com/labstack/echo/v4"
	apiMiddleware "github.com/speedianet/sam/src/presentation/api/middleware"
	"github.com/speedianet/sam/src/presentation/shared"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

// @title			SamApi
// @version			0.0.1
// @description		Speedia AppManager API
// @termsOfService	https://speedia.net/tos/

// @contact.name	Speedia Engineering
// @contact.url		https://speedia.net/
// @contact.email	eng+swagger@speedia.net

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey	Bearer
// @in 							header
// @name						Authorization
// @description					Type "Bearer" + JWT token or API key.

// @host		localhost:10000
// @BasePath	/v1
func ApiInit() {
	shared.CheckEnvs()

	e := echo.New()

	basePath := "/v1"
	baseRoute := e.Group(basePath)

	e.Pre(apiMiddleware.TrailingSlash(basePath))
	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders)
	e.Use(apiMiddleware.Auth(basePath))

	registerApiRoutes(baseRoute)

	e.Start(":10000")
}
