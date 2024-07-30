package cliController

import (
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type O11yController struct {
	o11yService *service.O11yService
}

func NewO11yController(
	transientDbService *internalDbInfra.TransientDatabaseService,
) *O11yController {
	return &O11yController{
		o11yService: service.NewO11yService(transientDbService),
	}
}

func (controller *O11yController) ReadOverview() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "GetO11yOverview",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.o11yService.ReadOverview())
		},
	}

	return cmd
}
