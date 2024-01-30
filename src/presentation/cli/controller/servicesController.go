package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetServicesController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetServices",
		Run: func(cmd *cobra.Command, args []string) {
			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesList, err := useCase.GetServicesWithMetrics(servicesQueryRepo)
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
			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesList, err := useCase.GetInstallableServices(servicesQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, servicesList)
		},
	}

	return cmd
}

func AddInstallableServiceController() *cobra.Command {
	var nameStr string
	var versionStr string
	var startupFileStr string
	var portBindingsSlice []string
	var autoCreateMapping bool

	cmd := &cobra.Command{
		Use:   "add-installable",
		Short: "AddInstallableService",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic(nameStr)

			var svcVersionPtr *valueObject.ServiceVersion
			if versionStr != "" {
				svcVersion := valueObject.NewServiceVersionPanic(versionStr)
				svcVersionPtr = &svcVersion
			}

			var startupFilePtr *valueObject.UnixFilePath
			if startupFileStr != "" {
				startupFile := valueObject.NewUnixFilePathPanic(startupFileStr)
				startupFilePtr = &startupFile
			}

			var portBindings []valueObject.PortBinding
			for _, portBinding := range portBindingsSlice {
				svcPortBinding, err := valueObject.NewPortBindingFromString(portBinding)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				portBindings = append(portBindings, svcPortBinding)
			}

			addInstallableServiceDto := dto.NewAddInstallableService(
				svcName,
				svcVersionPtr,
				startupFilePtr,
				portBindings,
				autoCreateMapping,
			)

			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

			err := useCase.AddInstallableService(
				servicesQueryRepo,
				servicesCmdRepo,
				vhostQueryRepo,
				vhostCmdRepo,
				addInstallableServiceDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "InstallableServiceAdded")
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&versionStr, "version", "v", "", "ServiceVersion")
	cmd.Flags().StringVarP(&startupFileStr, "startup-file", "f", "", "StartupFile")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "p", []string{}, "PortBindings (port/protocol)",
	)
	cmd.Flags().BoolVarP(
		&autoCreateMapping,
		"auto-create-mapping",
		"a",
		true,
		"AutoCreateMapping",
	)
	return cmd
}

func AddCustomServiceController() *cobra.Command {
	var nameStr string
	var typeStr string
	var commandStr string
	var versionStr string
	var portBindingsSlice []string
	var autoCreateMapping bool

	cmd := &cobra.Command{
		Use:   "add-custom",
		Short: "AddCustomService",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic(nameStr)
			svcType := valueObject.NewServiceTypePanic(typeStr)
			svcCommand := valueObject.NewUnixCommandPanic(commandStr)

			var svcVersionPtr *valueObject.ServiceVersion
			if versionStr != "" {
				svcVersion := valueObject.NewServiceVersionPanic(versionStr)
				svcVersionPtr = &svcVersion
			}

			var portBindings []valueObject.PortBinding
			for _, portBinding := range portBindingsSlice {
				svcPortBinding, err := valueObject.NewPortBindingFromString(portBinding)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				portBindings = append(portBindings, svcPortBinding)
			}

			addCustomServiceDto := dto.NewAddCustomService(
				svcName,
				svcType,
				svcCommand,
				svcVersionPtr,
				portBindings,
				autoCreateMapping,
			)

			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

			err := useCase.AddCustomService(
				servicesQueryRepo,
				servicesCmdRepo,
				vhostQueryRepo,
				vhostCmdRepo,
				addCustomServiceDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "CustomServiceAdded")
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&typeStr, "type", "t", "", "ServiceType (application|database|runtime|other)")
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "UnixCommand")
	cmd.MarkFlagRequired("command")
	cmd.Flags().StringVarP(&versionStr, "version", "v", "", "ServiceVersion")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "p", []string{}, "PortBindings (port/protocol)",
	)
	cmd.Flags().BoolVarP(
		&autoCreateMapping,
		"auto-create-mapping",
		"a",
		true,
		"AutoCreateMapping",
	)
	return cmd
}

func UpdateServiceController() *cobra.Command {
	var nameStr string
	var typeStr string
	var commandStr string
	var statusStr string
	var versionStr string
	var startupFileStr string
	var portBindingsSlice []string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateService",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic(nameStr)

			var svcTypePtr *valueObject.ServiceType
			if typeStr != "" {
				svcType := valueObject.NewServiceTypePanic(typeStr)
				svcTypePtr = &svcType
			}

			var svcStatusPtr *valueObject.ServiceStatus
			if statusStr != "" {
				svcStatus := valueObject.NewServiceStatusPanic(statusStr)
				svcStatusPtr = &svcStatus
			}

			var svcCommandPtr *valueObject.UnixCommand
			if commandStr != "" {
				svcCommand := valueObject.NewUnixCommandPanic(commandStr)
				svcCommandPtr = &svcCommand
			}

			var svcVersionPtr *valueObject.ServiceVersion
			if versionStr != "" {
				svcVersion := valueObject.NewServiceVersionPanic(versionStr)
				svcVersionPtr = &svcVersion
			}

			var svcStartupFilePtr *valueObject.UnixFilePath
			if startupFileStr != "" {
				svcStartupFile := valueObject.NewUnixFilePathPanic(startupFileStr)
				svcStartupFilePtr = &svcStartupFile
			}

			var portBindings []valueObject.PortBinding
			for _, portBinding := range portBindingsSlice {
				svcPortBinding, err := valueObject.NewPortBindingFromString(portBinding)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				portBindings = append(portBindings, svcPortBinding)
			}

			updateSvcDto := dto.NewUpdateService(
				svcName,
				svcTypePtr,
				svcCommandPtr,
				svcStatusPtr,
				svcVersionPtr,
				svcStartupFilePtr,
				portBindings,
			)

			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
			vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

			err := useCase.UpdateService(
				servicesQueryRepo,
				servicesCmdRepo,
				vhostQueryRepo,
				vhostCmdRepo,
				updateSvcDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ServiceUpdated")
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&typeStr, "type", "t", "", "ServiceType")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "UnixCommand")
	cmd.Flags().StringVarP(&statusStr, "status", "s", "", "ServiceStatus")
	cmd.Flags().StringVarP(&versionStr, "version", "v", "", "ServiceVersion")
	cmd.Flags().StringVarP(&startupFileStr, "startup-file", "f", "", "StartupFile")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "p", []string{}, "PortBindings (port/protocol)",
	)
	return cmd
}

func DeleteServiceController() *cobra.Command {
	var nameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteService",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic(nameStr)

			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesCmdRepo := servicesInfra.ServicesCmdRepo{}

			err := useCase.DeleteService(
				servicesQueryRepo,
				servicesCmdRepo,
				svcName,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "ServiceDeleted")
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	return cmd
}
