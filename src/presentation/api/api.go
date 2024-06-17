package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiMiddleware "github.com/speedianet/os/src/presentation/api/middleware"
	sharedMiddleware "github.com/speedianet/os/src/presentation/shared/middleware"
)

// @title			OsApi
// @version			0.0.2
// @description		Speedia OS API
// @termsOfService	https://speedia.net/tos/

// @contact.name	Speedia Engineering
// @contact.url		https://speedia.net/
// @contact.email	eng+swagger@speedia.net

// @license.name  Eclipse Public License v2.0
// @license.url   https://www.eclipse.org/legal/epl-2.0/

// @securityDefinitions.apikey	Bearer
// @in 							header
// @name						Authorization
// @description					Type "Bearer" + JWT token or API key.

// @host		localhost:1618
// @BasePath	/_/api
func ApiInit(
	e *echo.Echo,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) {
	sharedMiddleware.CheckEnvs()

	basePath := "/_/api"
	baseRoute := e.Group(basePath)

	e.Pre(apiMiddleware.AddTrailingSlash(basePath))

	requestTimeout := 180 * time.Second
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: requestTimeout,
	}))

	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders(basePath))
	e.Use(apiMiddleware.ReadOnlyMode(basePath))
	e.Use(apiMiddleware.Auth(basePath))

	router := NewRouter(baseRoute, transientDbSvc, persistentDbSvc)
	router.RegisterRoutes()
}
