package api

import (
	"time"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiMiddleware "github.com/goinfinite/os/src/presentation/api/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	ApiBasePath string = "/api"
)

// @title			OsApi
// @version			0.2.6
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
	e *echo.Echo,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) {
	baseRoute := e.Group(ApiBasePath)

	e.Pre(apiMiddleware.AddTrailingSlash(ApiBasePath))

	requestTimeout := 180 * time.Second
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: requestTimeout,
	}))

	e.Use(apiMiddleware.PanicHandler)
	e.Use(apiMiddleware.SetDefaultHeaders(ApiBasePath))
	e.Use(apiMiddleware.ReadOnlyMode(ApiBasePath))
	e.Use(apiMiddleware.SetDatabaseServices(
		persistentDbSvc, transientDbSvc, trailDbSvc,
	))
	e.Use(apiMiddleware.Authentication(ApiBasePath))

	router := NewRouter(baseRoute, transientDbSvc, persistentDbSvc, trailDbSvc)
	router.RegisterRoutes()
}
