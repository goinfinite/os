package cli

import (
	"fmt"
	"os"
	"path/filepath"

	cliInit "github.com/goinfinite/os/src/presentation/cli/init"
	cliMiddleware "github.com/goinfinite/os/src/presentation/cli/middleware"
	tkPresentationMiddleware "github.com/goinfinite/tk/src/presentation/middleware"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Infinite OS CLI",
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
	defer tkPresentationMiddleware.CliPanicHandler()
	cliMiddleware.PreventRootless()

	cliMiddleware.CheckEnvs()
	tkPresentationMiddleware.LogHandler{}.Init()

	transientDbSvc := cliInit.TransientDatabaseService()
	persistentDbSvc := cliInit.PersistentDatabaseService()
	trailDbSvc := cliInit.TrailDatabaseService()

	cliMiddleware.NewPrimaryVirtualHostSynchronizer(persistentDbSvc).Run()

	router := NewRouter(transientDbSvc, persistentDbSvc, trailDbSvc)
	router.RegisterRoutes()

	RunRootCmd()
}
