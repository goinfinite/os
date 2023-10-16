package cliMiddleware

import (
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/speedianet/sam/src/presentation/shared"
	"github.com/spf13/cobra"
)

func ServiceStatusValidator(serviceNameStr string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := shared.CheckServices(serviceNameStr)
		if err != nil {
			cliHelper.ResponseWrapper(false, err.Error())
		}
	}
}
