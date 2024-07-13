package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiMiddleware "github.com/speedianet/os/src/presentation/api/middleware"
	sharedMiddleware "github.com/speedianet/os/src/presentation/shared/middleware"
)

const (
	ApiBasePath string = "/api"
)

// @title			OsApi
// @version			0.0.4
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
// @BasePath	/api
func ApiInit(
	e *echo.Echo,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) {
	sharedMiddleware.CheckEnvs()

	baseRoute := e.Group(ApiBasePath)

	e.Pre(apiMiddleware.AddTrailingSlash(ApiBasePath))

	requestTimeout := 180 * time.Second
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: requestTimeout,
	}))

	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders(ApiBasePath))
	e.Use(apiMiddleware.ReadOnlyMode(ApiBasePath))
	e.Use(apiMiddleware.Auth(ApiBasePath))

	router := NewRouter(baseRoute, transientDbSvc, persistentDbSvc)
	router.RegisterRoutes()
}
