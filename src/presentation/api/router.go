package api

import (
	_ "embed"
	"net/http"
	"net/url"
	"strings"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiController "github.com/goinfinite/os/src/presentation/api/controller"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/goinfinite/os/src/presentation/api/docs"
)

type Router struct {
	baseRoute       *echo.Group
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewRouter(
	baseRoute *echo.Group,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *Router {
	return &Router{
		baseRoute:       baseRoute,
		transientDbSvc:  transientDbSvc,
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (router Router) swaggerRoute() {
	swaggerGroup := router.baseRoute.Group("/swagger")
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)
}

func (router Router) authRoutes() {
	authGroup := router.baseRoute.Group("/v1/auth")
	authController := apiController.NewAuthenticationController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	authGroup.POST("/login/", authController.Login)
}

func (router Router) accountRoutes() {
	accountGroup := router.baseRoute.Group("/v1/account")
	accountController := apiController.NewAccountController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	accountGroup.GET("/", accountController.Read)
	accountGroup.POST("/", accountController.Create)
	accountGroup.PUT("/", accountController.Update)
	accountGroup.DELETE("/:accountId/", accountController.Delete)

	secureAccessPublicKeyGroup := accountGroup.Group("/secure-access-public-key")
	secureAccessPublicKeyGroup.POST("/", accountController.CreateSecureAccessPublicKey)
	secureAccessPublicKeyGroup.DELETE(
		"/:secureAccessPublicKeyId/", accountController.DeleteSecureAccessPublicKey,
	)
}

func (router Router) cronRoutes() {
	cronGroup := router.baseRoute.Group("/v1/cron")
	cronController := apiController.NewCronController(router.trailDbSvc)

	cronGroup.GET("/", cronController.Read)
	cronGroup.POST("/", cronController.Create)
	cronGroup.PUT("/", cronController.Update)
	cronGroup.DELETE("/:cronId/", cronController.Delete)
}

func (router Router) databaseRoutes() {
	databaseGroup := router.baseRoute.Group("/v1/database")
	databaseController := apiController.NewDatabaseController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	databaseGroup.GET("/:dbType/", databaseController.Read)
	databaseGroup.POST("/:dbType/", databaseController.Create)
	databaseGroup.DELETE("/:dbType/:dbName/", databaseController.Delete)
	databaseGroup.POST("/:dbType/:dbName/user/", databaseController.CreateUser)
	databaseGroup.POST("/:dbType/user/", databaseController.CreateUser)
	databaseGroup.DELETE(
		"/:dbType/:dbName/user/:dbUser/", databaseController.DeleteUser,
	)
}

func (router Router) filesRoutes() {
	filesGroup := router.baseRoute.Group("/v1/files")

	filesController := apiController.NewFilesController(router.trailDbSvc)

	filesGroup.GET("/", filesController.Read)
	filesGroup.POST("/", filesController.Create)
	filesGroup.PUT("/", filesController.Update)
	filesGroup.POST("/copy/", filesController.Copy)
	filesGroup.PUT("/delete/", filesController.Delete)
	filesGroup.POST("/compress/", filesController.Compress)
	filesGroup.PUT("/extract/", filesController.Extract)
	filesGroup.POST("/upload/", filesController.Upload)
	filesGroup.GET("/download/", filesController.Download)
}

func (router Router) marketplaceRoutes() {
	marketplaceGroup := router.baseRoute.Group("/v1/marketplace")
	marketplaceController := apiController.NewMarketplaceController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	marketplaceInstalledGroup := marketplaceGroup.Group("/installed")
	marketplaceInstalledGroup.GET("/", marketplaceController.ReadInstalledItems)
	marketplaceInstalledGroup.DELETE(
		"/:installedId/",
		marketplaceController.DeleteInstalledItem,
	)

	marketplaceCatalogGroup := marketplaceGroup.Group("/catalog")
	marketplaceCatalogGroup.GET("/", marketplaceController.ReadCatalog)
	marketplaceCatalogGroup.POST("/", marketplaceController.InstallCatalogItem)

	go marketplaceController.AutoRefreshMarketplaceCatalogItems()
}

func (router Router) o11yRoutes() {
	o11yGroup := router.baseRoute.Group("/v1/o11y")

	o11yController := apiController.NewO11yController(router.transientDbSvc)
	o11yGroup.GET("/overview/", o11yController.ReadOverview)
}

func (router Router) runtimeRoutes() {
	runtimeGroup := router.baseRoute.Group("/v1/runtime")
	runtimeController := apiController.NewRuntimeController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	runtimeGroup.GET("/php/:hostname/", runtimeController.ReadPhpConfigs)
	runtimeGroup.PUT("/php/:hostname/", runtimeController.UpdatePhpConfigs)
}

func (router *Router) scheduledTaskRoutes() {
	scheduledTaskGroup := router.baseRoute.Group("/v1/scheduled-task")

	scheduledTaskController := apiController.NewScheduledTaskController(router.persistentDbSvc)
	scheduledTaskGroup.GET("/", scheduledTaskController.Read)
	scheduledTaskGroup.PUT("/", scheduledTaskController.Update)
	go scheduledTaskController.Run()
}

func (router Router) servicesRoutes() {
	servicesGroup := router.baseRoute.Group("/v1/services")
	servicesController := apiController.NewServicesController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	servicesGroup.GET("/", servicesController.ReadInstalledItems)
	servicesGroup.GET("/installables/", servicesController.ReadInstallablesItems)
	servicesGroup.POST("/installables/", servicesController.CreateInstallable)
	servicesGroup.POST("/custom/", servicesController.CreateCustom)
	servicesGroup.PUT("/", servicesController.Update)
	servicesGroup.DELETE("/:svcName/", servicesController.Delete)

	go servicesController.AutoRefreshServiceInstallableItems()
}

func (router Router) setupRoutes() {
	setupGroup := router.baseRoute.Group("/v1/setup")
	setupController := apiController.NewSetupController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	setupGroup.POST("/", setupController.Setup)
}

func (router Router) sslRoutes() {
	sslGroup := router.baseRoute.Group("/v1/ssl")
	sslController := apiController.NewSslController(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)

	sslGroup.GET("/", sslController.Read)
	sslGroup.POST("/", sslController.Create)
	sslGroup.POST("/trusted/", sslController.CreatePubliclyTrusted)
	sslGroup.DELETE("/:sslPairId/", sslController.Delete)

	go sslController.SslCertificateWatchdog()
}

func (router Router) vhostRoutes() {
	vhostsGroup := router.baseRoute.Group("/v1/vhosts")
	vhostsGroup.Any("/*", func(c echo.Context) error {
		originalPath := c.Request().URL.Path
		parsedPath, err := url.Parse(originalPath)
		if err != nil {
			return c.String(http.StatusBadRequest, "InvalidUrl")
		}
		newPath := strings.ReplaceAll(parsedPath.Path, "/v1/vhosts", "/v1/vhost")
		return c.Redirect(http.StatusTemporaryRedirect, newPath)
	})

	vhostGroup := router.baseRoute.Group("/v1/vhost")
	vhostController := apiController.NewVirtualHostController(
		router.persistentDbSvc, router.trailDbSvc,
	)

	vhostGroup.GET("/", vhostController.Read)
	vhostGroup.POST("/", vhostController.Create)
	vhostGroup.PUT("/", vhostController.Update)
	vhostGroup.DELETE("/:hostname/", vhostController.Delete)

	mappingsGroup := vhostGroup.Group("/mapping")
	mappingsGroup.GET("/", vhostController.ReadWithMappings)
	mappingsGroup.POST("/", vhostController.CreateMapping)
	mappingsGroup.PUT("/", vhostController.UpdateMapping)
	mappingsGroup.DELETE(
		"/:mappingId/",
		vhostController.DeleteMapping,
	)

	mappingSecurityRuleGroup := mappingsGroup.Group("/security-rule")
	mappingSecurityRuleGroup.GET("/", vhostController.ReadMappingSecurityRules)
	mappingSecurityRuleGroup.POST("/", vhostController.CreateMappingSecurityRule)
	mappingSecurityRuleGroup.PUT("/:id/", vhostController.UpdateMappingSecurityRule)
	mappingSecurityRuleGroup.DELETE("/:id/", vhostController.DeleteMappingSecurityRule)
}

func (router Router) RegisterRoutes() {
	router.swaggerRoute()
	router.authRoutes()
	router.accountRoutes()
	router.cronRoutes()
	router.databaseRoutes()
	router.filesRoutes()
	router.marketplaceRoutes()
	router.o11yRoutes()
	router.runtimeRoutes()
	router.scheduledTaskRoutes()
	router.servicesRoutes()
	router.setupRoutes()
	router.sslRoutes()
	router.vhostRoutes()
}
