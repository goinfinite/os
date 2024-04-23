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
	marketplaceInstalledsGroup.GET("/", marketplaceController.GetInstalledItems)
	marketplaceInstalledsGroup.DELETE(
		"/:installedId/",
		marketplaceController.DeleteInstalledItem,
	)

	marketplaceCatalogGroup := marketplaceGroup.Group("/catalog")
	marketplaceCatalogGroup.GET("/", marketplaceController.GetCatalog)
	marketplaceCatalogGroup.POST("/", marketplaceController.InstallCatalogItem)
}

func (router Router) o11yRoutes(baseRoute *echo.Group) {
	o11yGroup := baseRoute.Group("/o11y")

	o11yController := apiController.NewO11yController(router.transientDbSvc)
	o11yGroup.GET("/overview/", o11yController.GetO11yOverview)
}

func (router Router) runtimeRoutes(baseRoute *echo.Group) {
	runtimeGroup := baseRoute.Group("/runtime")
	runtimeGroup.GET("/php/:hostname/", apiController.GetPhpConfigsController)
	runtimeGroup.PUT("/php/:hostname/", apiController.UpdatePhpConfigsController)
}

func (router Router) servicesRoutes(baseRoute *echo.Group) {
	servicesGroup := baseRoute.Group("/services")
	servicesGroup.GET("/", apiController.GetServicesController)
	servicesGroup.GET("/installables/", apiController.GetInstallableServicesController)
	servicesGroup.POST("/installables/", apiController.CreateInstallableServiceController)
	servicesGroup.POST("/custom/", apiController.CreateCustomServiceController)
	servicesGroup.PUT("/", apiController.UpdateServiceController)
	servicesGroup.DELETE("/:svcName/", apiController.DeleteServiceController)
}

func (router Router) sslRoutes(baseRoute *echo.Group) {
	sslGroup := baseRoute.Group("/ssl")
	sslGroup.GET("/", apiController.GetSslPairsController)
	sslGroup.POST("/", apiController.CreateSslPairController)
	sslGroup.DELETE("/:sslPairId/", apiController.DeleteSslPairController)
	go apiController.SslCertificateWatchdogController()
}

func (router Router) vhostsRoutes(baseRoute *echo.Group) {
	vhostsGroup := baseRoute.Group("/vhosts")
	vhostsGroup.GET("/", apiController.GetVirtualHostsController)
	vhostsGroup.POST("/", apiController.CreateVirtualHostController)
	vhostsGroup.DELETE("/:hostname/", apiController.DeleteVirtualHostController)

	vhostsGroup.GET("/mapping/", apiController.GetVirtualHostsWithMappingsController)
	vhostsGroup.POST("/mapping/", apiController.CreateVirtualHostMappingController)
	vhostsGroup.DELETE(
		"/mapping/:hostname/:mappingId/",
		apiController.DeleteVirtualHostMappingController,
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
