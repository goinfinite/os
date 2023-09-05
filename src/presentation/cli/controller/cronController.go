package cliController

import (
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetCronsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetCrons",
		Run: func(cmd *cobra.Command, args []string) {
			cronQueryRepo := infra.CronQueryRepo{}

			cronsList, err := useCase.GetCrons(cronQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, cronsList)
		},
	}
	return cmd
}
