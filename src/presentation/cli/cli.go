package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cliInit "github.com/speedianet/os/src/presentation/cli/init"
	cliMiddleware "github.com/speedianet/os/src/presentation/cli/middleware"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Speedia OS CLI",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func RunRootCmd() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func CliInit() {
	defer cliMiddleware.PanicHandler()
	cliMiddleware.PreventRootless()

	cliMiddleware.CheckEnvs()
	cliMiddleware.LogHandler()

	transientDbSvc := cliInit.TransientDatabaseService()
	persistentDbSvc := cliInit.PersistentDatabaseService()

	router := NewRouter(transientDbSvc, persistentDbSvc)
	router.RegisterRoutes()

	RunRootCmd()
}
