package cliController

import (
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetSslsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetSsls",
		Run: func(cmd *cobra.Command, args []string) {
			sslQueryRepo := infra.NewSslQueryRepo()
			sslsList, err := useCase.GetSsls(sslQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, sslsList)
		},
	}

	return cmd
}
