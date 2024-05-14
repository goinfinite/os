package cli

import (
	"fmt"

	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	api "github.com/speedianet/os/src/presentation/api"
	cliController "github.com/speedianet/os/src/presentation/cli/controller"
	"github.com/spf13/cobra"
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

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "ShowSoftwareVersion",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Speedia OS v0.0.1")
	},
}

func (router Router) accountRoutes() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "AccountManagement",
	}

	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(cliController.GetAccountsController())
	accountCmd.AddCommand(cliController.CreateAccountController())
	accountCmd.AddCommand(cliController.DeleteAccountController())
	accountCmd.AddCommand(cliController.UpdateAccountController())
}

func (router Router) cronRoutes() {
	var cronCmd = &cobra.Command{
		Use:   "cron",
		Short: "CronManagement",
	}

	rootCmd.AddCommand(cronCmd)
	cronCmd.AddCommand(cliController.GetCronsController())
	cronCmd.AddCommand(cliController.CreateCronControler())
	cronCmd.AddCommand(cliController.UpdateCronController())
	cronCmd.AddCommand(cliController.DeleteCronController())
}

func (router Router) databaseRoutes() {
	var databaseCmd = &cobra.Command{
		Use:   "db",
		Short: "DatabaseManagement",
	}

	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(cliController.GetDatabasesController())
	databaseCmd.AddCommand(cliController.CreateDatabaseController())
	databaseCmd.AddCommand(cliController.DeleteDatabaseController())
	databaseCmd.AddCommand(cliController.CreateDatabaseUserController())
	databaseCmd.AddCommand(cliController.DeleteDatabaseUserController())
}

func (router Router) marketplaceRoutes() {
	var marketplaceCmd = &cobra.Command{
		Use:   "mktplace",
		Short: "Marketplace",
	}

	rootCmd.AddCommand(marketplaceCmd)

	marketplaceController := cliController.NewMarketplaceController(
		router.persistentDbSvc,
	)
	marketplaceCmd.AddCommand(marketplaceController.GetCatalog())
	marketplaceCmd.AddCommand(marketplaceController.InstallCatalogItem())
}

func (router Router) o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	rootCmd.AddCommand(o11yCmd)
	o11yCmd.AddCommand(cliController.ReadO11yOverviewController(router.transientDbSvc))
}

func (router Router) runtimeRoutes() {
	var runtimeCmd = &cobra.Command{
		Use:   "runtime",
		Short: "RuntimeManagement",
	}

	var phpCmd = &cobra.Command{
		Use:   "php",
		Short: "PhpManagement",
	}

	rootCmd.AddCommand(runtimeCmd)
	runtimeCmd.AddCommand(phpCmd)
	phpCmd.AddCommand(cliController.GetPhpConfigsController())
	phpCmd.AddCommand(cliController.UpdatePhpConfigController())
	phpCmd.AddCommand(cliController.UpdatePhpSettingController())
	phpCmd.AddCommand(cliController.UpdatePhpModuleController())
}

func (router Router) serveRoutes() {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start Speedia OS HTTPS server (port 1618)",
		Run: func(cmd *cobra.Command, args []string) {
			api.ApiInit(router.transientDbSvc, router.persistentDbSvc)
		},
	}

	rootCmd.AddCommand(serveCmd)
}

func (router Router) servicesRoutes() {
	var servicesCmd = &cobra.Command{
		Use:   "services",
		Short: "ServicesManagement",
	}

	rootCmd.AddCommand(servicesCmd)

	servicesController := cliController.NewServicesController(
		router.persistentDbSvc,
	)
	servicesCmd.AddCommand(servicesController.Read())
	servicesCmd.AddCommand(servicesController.ReadInstallables())
	servicesCmd.AddCommand(servicesController.CreateInstallable())
	servicesCmd.AddCommand(servicesController.CreateCustom())
	servicesCmd.AddCommand(servicesController.Update())
	servicesCmd.AddCommand(servicesController.Delete())
}

func (router Router) sslRoutes() {
	var sslCmd = &cobra.Command{
		Use:   "ssl",
		Short: "SslManagement",
	}

	rootCmd.AddCommand(sslCmd)

	sslController := cliController.NewSslController(
		router.persistentDbSvc,
	)
	sslCmd.AddCommand(sslController.Read())
	sslCmd.AddCommand(sslController.Create())
	sslCmd.AddCommand(sslController.DeleteVhosts())
	sslCmd.AddCommand(sslController.Delete())
}

func (router Router) virtualHostRoutes() {
	var vhostCmd = &cobra.Command{
		Use:   "vhost",
		Short: "VirtualHostManagement",
	}

	rootCmd.AddCommand(vhostCmd)

	vhostController := cliController.NewVirtualHostController(
		router.persistentDbSvc,
	)
	vhostCmd.AddCommand(vhostController.Get())
	vhostCmd.AddCommand(vhostController.Create())
	vhostCmd.AddCommand(vhostController.Delete())

	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}

	vhostCmd.AddCommand(mappingCmd)
	mappingCmd.AddCommand(vhostController.GetWithMappings())
	mappingCmd.AddCommand(vhostController.CreateMapping())
	mappingCmd.AddCommand(vhostController.DeleteMapping())
}

func (router Router) RegisterRoutes() {
	rootCmd.AddCommand(versionCmd)

	router.accountRoutes()
	router.cronRoutes()
	router.databaseRoutes()
	router.marketplaceRoutes()
	router.o11yRoutes()
	router.runtimeRoutes()
	router.serveRoutes()
	router.servicesRoutes()
	router.sslRoutes()
	router.virtualHostRoutes()
}
