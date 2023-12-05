package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetServicesController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetServices",
		Run: func(cmd *cobra.Command, args []string) {
			servicesQueryRepo := infra.ServicesQueryRepo{}
			servicesList, err := useCase.GetServices(servicesQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, servicesList)
		},
	}

	return cmd
}

func GetInstallableServicesController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-installables",
		Short: "GetInstallableServices",
		Run: func(cmd *cobra.Command, args []string) {
			servicesQueryRepo := infra.ServicesQueryRepo{}
			servicesList, err := useCase.GetInstallableServices(servicesQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, servicesList)
		},
	}

	return cmd
}

func UpdateServiceController() *cobra.Command {
	var nameStr string
	var statusStr string
	var versionStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateServiceStatus",
		Run: func(cmd *cobra.Command, args []string) {

			svcName := valueObject.NewServiceNamePanic(nameStr)
			svcStatus := valueObject.NewServiceStatusPanic(statusStr)
			var svcVersionPtr *valueObject.ServiceVersion
			if versionStr != "" {
				svcVersion := valueObject.NewServiceVersionPanic(versionStr)
				svcVersionPtr = &svcVersion
			}

			updateSvcStatusDto := dto.NewUpdateSvcStatus(
				svcName,
				svcStatus,
				svcVersionPtr,
			)

			servicesQueryRepo := infra.ServicesQueryRepo{}
			servicesCmdRepo := infra.ServicesCmdRepo{}

			err := useCase.UpdateServiceStatus(
				servicesQueryRepo,
				servicesCmdRepo,
				updateSvcStatusDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ServiceStatusUpdated")
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&statusStr, "status", "s", "", "ServiceStatus")
	cmd.MarkFlagRequired("status")
	cmd.Flags().StringVarP(&versionStr, "version", "v", "", "ServiceVersion")
	return cmd
}
