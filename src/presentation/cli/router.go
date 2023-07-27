package cli

import (
	"fmt"

	api "github.com/speedianet/sam/src/presentation/api"
	cliController "github.com/speedianet/sam/src/presentation/cli/controller"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Speedia AppManager v0.0.1")
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SAM server (default to port 10000)",
	Run: func(cmd *cobra.Command, args []string) {
		api.ApiInit()
	},
}

func databaseRoutes() {
	var databaseCmd = &cobra.Command{
		Use:   "db",
		Short: "DatabaseManagement",
	}

	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(cliController.GetDatabasesController())
	databaseCmd.AddCommand(cliController.AddDatabaseController())
	databaseCmd.AddCommand(cliController.DeleteDatabaseController())
	databaseCmd.AddCommand(cliController.AddDatabaseUserController())
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

func userRoutes() {
	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "UserManagement",
	}

	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(cliController.GetUsersController())
	userCmd.AddCommand(cliController.AddUserController())
	userCmd.AddCommand(cliController.DeleteUserController())
	userCmd.AddCommand(cliController.UpdateUserController())
}

func registerCliRoutes() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	databaseRoutes()
	o11yRoutes()
	userRoutes()
}
