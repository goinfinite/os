package cliMiddleware

import (
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func ServiceStatusValidator(serviceNameStr string) func(cmd *cobra.Command, args []string) {
	servicesQueryRepo := infra.ServicesQueryRepo{}

	return func(cmd *cobra.Command, args []string) {
		serviceName, err := valueObject.NewServiceName(serviceNameStr)

		currentSvcStatus, err := servicesQueryRepo.GetByName(serviceName)
		if err != nil {
			cliHelper.ResponseWrapper(false, "Failed to get "+serviceNameStr+" service: "+err.Error())
		}

		var badCommandMessage string

		isStopped := currentSvcStatus.Status.String() == "stopped"
		if isStopped {
			badCommandMessage = "Service paused"
		}
		isUninstalled := currentSvcStatus.Status.String() == "uninstalled"
		if isUninstalled {
			badCommandMessage = "Service not installed"
		}
		shouldInstall := isStopped || isUninstalled
		if shouldInstall {
			cliHelper.ResponseWrapper(false, badCommandMessage)
		}
	}
}
