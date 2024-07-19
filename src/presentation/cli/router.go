package cli

import (
	"fmt"

	infraEnvs "github.com/speedianet/os/src/infra/envs"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"github.com/speedianet/os/src/presentation"
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
		fmt.Println("Speedia OS v" + infraEnvs.SpeediaOsVersion)
	},
}

func (router Router) accountRoutes() {
	var accountCmd = &cobra.Command{
		Use:   "account",
		Short: "AccountManagement",
	}
	rootCmd.AddCommand(accountCmd)

	accountController := cliController.NewAccountController()

	accountCmd.AddCommand(accountController.Read())
	accountCmd.AddCommand(accountController.Create())
	accountCmd.AddCommand(accountController.Update())
	accountCmd.AddCommand(accountController.Delete())
}

func (router Router) authenticationRoutes() {
	var authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Authentication&Authorization",
	}
	rootCmd.AddCommand(authCmd)

	authenticationController := cliController.AuthController{}
	authCmd.AddCommand(authenticationController.Login())
}

func (router Router) cronRoutes() {
	var cronCmd = &cobra.Command{
		Use:   "cron",
		Short: "CronManagement",
	}
	rootCmd.AddCommand(cronCmd)

	cronController := cliController.NewCronController()
	cronCmd.AddCommand(cronController.Read())
	cronCmd.AddCommand(cronController.Create())
	cronCmd.AddCommand(cronController.Update())
	cronCmd.AddCommand(cronController.Delete())
}

func (router Router) databaseRoutes() {
	var databaseCmd = &cobra.Command{
		Use:   "db",
		Short: "DatabaseManagement",
	}
	rootCmd.AddCommand(databaseCmd)

	databaseController := cliController.NewDatabaseController(
		router.persistentDbSvc,
	)
	databaseCmd.AddCommand(databaseController.Read())
	databaseCmd.AddCommand(databaseController.Create())
	databaseCmd.AddCommand(databaseController.Delete())
	databaseCmd.AddCommand(databaseController.CreateUser())
	databaseCmd.AddCommand(databaseController.DeleteUser())
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
	marketplaceCmd.AddCommand(marketplaceController.ReadCatalog())
	marketplaceCmd.AddCommand(marketplaceController.InstallCatalogItem())
	marketplaceCmd.AddCommand(marketplaceController.ReadInstalledItems())
	marketplaceCmd.AddCommand(marketplaceController.DeleteInstalledItem())
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
	rootCmd.AddCommand(runtimeCmd)

	var phpCmd = &cobra.Command{
		Use:   "php",
		Short: "PhpManagement",
	}
	runtimeCmd.AddCommand(phpCmd)

	runtimeController := cliController.NewRuntimeController(router.persistentDbSvc)
	phpCmd.AddCommand(runtimeController.ReadPhpConfigs())
	phpCmd.AddCommand(runtimeController.UpdatePhpConfig())
	phpCmd.AddCommand(runtimeController.UpdatePhpSetting())
	phpCmd.AddCommand(runtimeController.UpdatePhpModule())
}

func (router *Router) scheduledTaskRoutes() {
	var scheduledTaskCmd = &cobra.Command{
		Use:   "task",
		Short: "ScheduledTaskManagement",
	}
	rootCmd.AddCommand(scheduledTaskCmd)

	scheduledTaskController := cliController.NewScheduledTaskController(router.persistentDbSvc)
	scheduledTaskCmd.AddCommand(scheduledTaskController.Read())
	scheduledTaskCmd.AddCommand(scheduledTaskController.Update())
}

func (router Router) serveRoutes() {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start Speedia OS HTTPS server (port 1618)",
		Run: func(cmd *cobra.Command, args []string) {
			presentation.HttpServerInit(router.persistentDbSvc, router.transientDbSvc)
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
		router.persistentDbSvc, router.transientDbSvc,
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
	vhostCmd.AddCommand(vhostController.Read())
	vhostCmd.AddCommand(vhostController.Create())
	vhostCmd.AddCommand(vhostController.Delete())

	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}
	vhostCmd.AddCommand(mappingCmd)

	mappingCmd.AddCommand(vhostController.ReadWithMappings())
	mappingCmd.AddCommand(vhostController.CreateMapping())
	mappingCmd.AddCommand(vhostController.DeleteMapping())
}

func (router Router) RegisterRoutes() {
	rootCmd.AddCommand(versionCmd)

	router.accountRoutes()
	router.authenticationRoutes()
	router.cronRoutes()
	router.databaseRoutes()
	router.marketplaceRoutes()
	router.o11yRoutes()
	router.runtimeRoutes()
	router.scheduledTaskRoutes()
	router.serveRoutes()
	router.servicesRoutes()
	router.sslRoutes()
	router.virtualHostRoutes()
}
