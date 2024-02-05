package cli

import (
	"fmt"

	api "github.com/speedianet/os/src/presentation/api"
	cliController "github.com/speedianet/os/src/presentation/cli/controller"
	cliMiddleware "github.com/speedianet/os/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Speedia OS v0.0.1")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SOS server (default to port 1618)",
	Run: func(cmd *cobra.Command, args []string) {
		api.ApiInit()
	},
}

func accountRoutes() {
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

func cronRoutes() {
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

func databaseRoutes() {
	var databaseCmd = &cobra.Command{
		Use:              "db",
		Short:            "DatabaseManagement",
		PersistentPreRun: cliMiddleware.ServiceStatusValidator("mysql"),
	}

	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(cliController.GetDatabasesController())
	databaseCmd.AddCommand(cliController.CreateDatabaseController())
	databaseCmd.AddCommand(cliController.DeleteDatabaseController())
	databaseCmd.AddCommand(cliController.CreateDatabaseUserController())
	databaseCmd.AddCommand(cliController.DeleteDatabaseUserController())
}

func o11yRoutes() {
	var o11yCmd = &cobra.Command{
		Use:   "o11y",
		Short: "O11yManagement",
	}

	rootCmd.AddCommand(o11yCmd)
	o11yCmd.AddCommand(cliController.GetO11yOverviewController())
}

func runtimeRoutes() {
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

func servicesRoutes() {
	var servicesCmd = &cobra.Command{
		Use:   "services",
		Short: "ServicesManagement",
	}

	rootCmd.AddCommand(servicesCmd)
	servicesCmd.AddCommand(cliController.GetServicesController())
	servicesCmd.AddCommand(cliController.GetInstallableServicesController())
	servicesCmd.AddCommand(cliController.CreateInstallableServiceController())
	servicesCmd.AddCommand(cliController.CreateCustomServiceController())
	servicesCmd.AddCommand(cliController.UpdateServiceController())
	servicesCmd.AddCommand(cliController.DeleteServiceController())
}

func sslRoutes() {
	var sslCmd = &cobra.Command{
		Use:   "ssl",
		Short: "SslManagement",
	}

	rootCmd.AddCommand(sslCmd)
	sslCmd.AddCommand(cliController.GetSslPairsController())
	sslCmd.AddCommand(cliController.CreateSslPairController())
	sslCmd.AddCommand(cliController.DeleteSslPairController())
}

func virtualHostRoutes() {
	var vhostCmd = &cobra.Command{
		Use:   "vhost",
		Short: "VirtualHostManagement",
	}

	rootCmd.AddCommand(vhostCmd)
	vhostCmd.AddCommand(cliController.GetVirtualHostsController())
	vhostCmd.AddCommand(cliController.AddVirtualHostController())
	vhostCmd.AddCommand(cliController.DeleteVirtualHostController())

	var mappingCmd = &cobra.Command{
		Use:   "mapping",
		Short: "MappingManagement",
	}

	vhostCmd.AddCommand(mappingCmd)
	mappingCmd.AddCommand(cliController.GetVirtualHostsWithMappingsController())
	mappingCmd.AddCommand(cliController.AddVirtualHostMappingController())
	mappingCmd.AddCommand(cliController.DeleteVirtualHostMappingController())
}

func registerCliRoutes() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	accountRoutes()
	cronRoutes()
	databaseRoutes()
	o11yRoutes()
	runtimeRoutes()
	servicesRoutes()
	sslRoutes()
	virtualHostRoutes()
}
