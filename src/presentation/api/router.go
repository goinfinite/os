package api

import (
	_ "embed"

	"github.com/labstack/echo/v4"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiController "github.com/speedianet/os/src/presentation/api/controller"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/speedianet/os/src/presentation/api/docs"
)

type Router struct {
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewRouter(
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *Router {
	return &Router{
		transientDbSvc:  transientDbSvc,
		persistentDbSvc: persistentDbSvc,
	}
}

func (router Router) swaggerRoute(baseRoute *echo.Group) {
	swaggerGroup := baseRoute.Group("/swagger")
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)
}

func (router Router) authRoutes(baseRoute *echo.Group) {
	authGroup := baseRoute.Group("/auth")
	authGroup.POST("/login/", apiController.AuthLoginController)
}

func (router Router) accountRoutes(baseRoute *echo.Group) {
	accountGroup := baseRoute.Group("/account")
	accountGroup.GET("/", apiController.GetAccountsController)
	accountGroup.POST("/", apiController.CreateAccountController)
	accountGroup.PUT("/", apiController.UpdateAccountController)
	accountGroup.DELETE("/:accountId/", apiController.DeleteAccountController)
}

func (router Router) cronRoutes(baseRoute *echo.Group) {
	cronGroup := baseRoute.Group("/cron")
	cronGroup.GET("/", apiController.GetCronsController)
	cronGroup.POST("/", apiController.CreateCronController)
	cronGroup.PUT("/", apiController.UpdateCronController)
	cronGroup.DELETE("/:cronId/", apiController.DeleteCronController)
}

func (router Router) databaseRoutes(baseRoute *echo.Group) {
	databaseGroup := baseRoute.Group("/database")
	databaseGroup.GET("/:dbType/", apiController.GetDatabasesController)
	databaseGroup.POST("/:dbType/", apiController.CreateDatabaseController)
	databaseGroup.DELETE(
		"/:dbType/:dbName/",
		apiController.DeleteDatabaseController,
	)
	databaseGroup.POST(
		"/:dbType/:dbName/user/",
		apiController.CreateDatabaseUserController,
	)
	databaseGroup.DELETE(
		"/:dbType/:dbName/user/:dbUser/",
		apiController.DeleteDatabaseUserController,
	)
}

func (router Router) filesRoutes(baseRoute *echo.Group) {
	filesGroup := baseRoute.Group("/files")
	filesGroup.GET("/", apiController.GetFilesController)
	filesGroup.POST("/", apiController.CreateFileController)
	filesGroup.PUT("/", apiController.UpdateFileController)
	filesGroup.POST("/copy/", apiController.CopyFileController)
	filesGroup.PUT("/delete/", apiController.DeleteFileController)
	filesGroup.POST("/compress/", apiController.CompressFilesController)
	filesGroup.PUT("/extract/", apiController.ExtractFilesController)
	filesGroup.POST("/upload/", apiController.UploadFilesController)
}

func (router Router) marketplaceRoutes(baseRoute *echo.Group) {
	marketplaceGroup := baseRoute.Group("/marketplace")
	marketplaceController := apiController.NewMarketplaceController(
		router.persistentDbSvc,
	)

	marketplaceInstalledsGroup := marketplaceGroup.Group("/installed")
	marketplaceInstalledsGroup.GET("/", marketplaceController.ReadInstalledItems)
	marketplaceInstalledsGroup.DELETE(
		"/:installedId/",
		marketplaceController.DeleteInstalledItem,
	)

	marketplaceCatalogGroup := marketplaceGroup.Group("/catalog")
	marketplaceCatalogGroup.GET("/", marketplaceController.ReadCatalog)
	marketplaceCatalogGroup.POST("/", marketplaceController.InstallCatalogItem)
}

func (router Router) o11yRoutes(baseRoute *echo.Group) {
	o11yGroup := baseRoute.Group("/o11y")

	o11yController := apiController.NewO11yController(router.transientDbSvc)
	o11yGroup.GET("/overview/", o11yController.ReadOverview)
}

func (router Router) runtimeRoutes(baseRoute *echo.Group) {
	runtimeGroup := baseRoute.Group("/runtime")
	runtimeController := apiController.NewRuntimeController(
		router.persistentDbSvc,
	)

	runtimeGroup.GET("/php/:hostname/", runtimeController.ReadConfigs)
	runtimeGroup.PUT("/php/:hostname/", runtimeController.UpdateConfigs)
}

func (router Router) servicesRoutes(baseRoute *echo.Group) {
	servicesGroup := baseRoute.Group("/services")
	servicesController := apiController.NewServicesController(
		router.persistentDbSvc,
	)

	servicesGroup.GET("/", servicesController.Read)
	servicesGroup.GET("/installables/", servicesController.ReadInstallables)
	servicesGroup.POST("/installables/", servicesController.CreateInstallable)
	servicesGroup.POST("/custom/", servicesController.CreateCustom)
	servicesGroup.PUT("/", servicesController.Update)
	servicesGroup.DELETE("/:svcName/", servicesController.Delete)
}

func (router Router) sslRoutes(baseRoute *echo.Group) {
	sslGroup := baseRoute.Group("/ssl")
	sslController := apiController.NewSslController(
		router.persistentDbSvc,
	)

	sslGroup.GET("/", sslController.Read)
	sslGroup.POST("/", sslController.Create)
	sslGroup.DELETE("/:sslPairId/", sslController.Delete)
	sslGroup.PUT("/vhost/", sslController.DeleteVhosts)
	go sslController.SslCertificateWatchdog()
}

func (router Router) vhostsRoutes(baseRoute *echo.Group) {
	vhostsGroup := baseRoute.Group("/vhosts")
	vhostController := apiController.NewVirtualHostController(
		router.persistentDbSvc,
	)

	vhostsGroup.GET("/", vhostController.Get)
	vhostsGroup.POST("/", vhostController.Create)
	vhostsGroup.DELETE("/:hostname/", vhostController.Delete)

	mappingsGroup := vhostsGroup.Group("/mapping")
	mappingsGroup.GET("/", vhostController.GetWithMappings)
	mappingsGroup.POST("/", vhostController.CreateMapping)
	mappingsGroup.DELETE(
		"/:mappingId/",
		vhostController.DeleteMapping,
	)
}

func (router Router) RegisterRoutes(baseRoute *echo.Group) {
	router.swaggerRoute(baseRoute)
	router.authRoutes(baseRoute)
	router.accountRoutes(baseRoute)
	router.cronRoutes(baseRoute)
	router.databaseRoutes(baseRoute)
	router.filesRoutes(baseRoute)
	router.marketplaceRoutes(baseRoute)
	router.o11yRoutes(baseRoute)
	router.runtimeRoutes(baseRoute)
	router.servicesRoutes(baseRoute)
	router.sslRoutes(baseRoute)
	router.vhostsRoutes(baseRoute)
}
