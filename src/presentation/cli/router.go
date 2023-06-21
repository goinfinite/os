package cli

import (
	"fmt"
	"os"
	"path/filepath"

	restApi "github.com/speedianet/sam/src/presentation/api"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Speedia AppManager CLI",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

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
		restApi.StartRestApi()
	},
}

func CliRouterInit() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serveCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
