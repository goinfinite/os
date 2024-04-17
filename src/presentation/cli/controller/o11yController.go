package cliController

import (
	"github.com/speedianet/os/src/domain/useCase"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	o11yInfra "github.com/speedianet/os/src/infra/o11y"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetO11yOverviewController(
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "GetOverview",
		Run: func(cmd *cobra.Command, args []string) {
			o11yQueryRepo := o11yInfra.NewO11yQueryRepo(transientDbSvc)
			o11yOverview, err := useCase.GetO11yOverview(o11yQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, o11yOverview)
		},
	}

	return cmd
}
