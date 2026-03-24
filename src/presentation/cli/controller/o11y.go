package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/spf13/cobra"
)

type O11yController struct {
	o11yLiaison *liaison.O11yLiaison
}

func NewO11yController(
	transientDbService *internalDbInfra.TransientDatabaseService,
) *O11yController {
	return &O11yController{
		o11yLiaison: liaison.NewO11yLiaison(transientDbService),
	}
}

func (controller *O11yController) ReadOverview() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "ReadO11yOverview",
		Run: func(cmd *cobra.Command, args []string) {
			tkPresentation.LiaisonCliResponseRenderer(controller.o11yLiaison.ReadOverview())
		},
	}

	return cmd
}
