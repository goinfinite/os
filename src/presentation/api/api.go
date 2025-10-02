package api

import (
	"time"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiMiddleware "github.com/goinfinite/os/src/presentation/api/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title			OsApi
// @version			0.2.7
// @description		Infinite OS API
// @termsOfService	https://goinfinite.net/tos/

// @contact.name	Infinite Engineering
// @contact.url		https://goinfinite.net/
// @contact.email	eng+swagger@goinfinite.net

// @license.name  Eclipse Public License v2.0
// @license.url   https://www.eclipse.org/legal/epl-2.0/

// @securityDefinitions.apikey	Bearer
// @in 							header
// @name						Authorization
// @description					Type "Bearer" + JWT token or API key.

// @host		localhost:1618
// @BasePath	/api
func ApiInit(
	echoInstance *echo.Echo,
	apiBasePath string,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) {
	baseRoute := echoInstance.Group(apiBasePath)

	echoInstance.Pre(apiMiddleware.AddTrailingSlash(apiBasePath))

	requestTimeout := 180 * time.Second
	echoInstance.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: requestTimeout,
	}))

	echoInstance.Use(apiMiddleware.SetDefaultHeaders(apiBasePath))
	echoInstance.Use(apiMiddleware.ReadOnlyMode(apiBasePath))
	echoInstance.Use(apiMiddleware.SetDatabaseServices(
		persistentDbSvc, transientDbSvc, trailDbSvc,
	))
	echoInstance.Use(apiMiddleware.Authentication(apiBasePath, persistentDbSvc))

	router := NewRouter(baseRoute, transientDbSvc, persistentDbSvc, trailDbSvc)
	router.RegisterRoutes()
}
