package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type ServicesController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServicesController {
	return &ServicesController{
		persistentDbSvc: persistentDbSvc,
	}
}

func (controller *ServicesController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
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

func (controller *ServicesController) ReadInstallables() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-installables",
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

func (controller *ServicesController) CreateInstallable() *cobra.Command {
	var nameStr string
	var versionStr string
	var startupFileStr string
	var portBindingsSlice []string
	var autoCreateMapping bool

	cmd := &cobra.Command{
		Use:   "create-installable",
		Short: "CreateInstallableService",
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

			createInstallableServiceDto := dto.NewCreateInstallableService(
				svcName,
				svcVersionPtr,
				startupFilePtr,
				portBindings,
				autoCreateMapping,
			)

			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

			err := useCase.CreateInstallableService(
				servicesQueryRepo,
				servicesCmdRepo,
				mappingQueryRepo,
				mappingCmdRepo,
				vhostQueryRepo,
				createInstallableServiceDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "InstallableServiceCreated")
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

func (controller *ServicesController) CreateCustom() *cobra.Command {
	var nameStr string
	var typeStr string
	var commandStr string
	var versionStr string
	var portBindingsSlice []string
	var autoCreateMapping bool

	cmd := &cobra.Command{
		Use:   "create-custom",
		Short: "CreateCustomService",
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

			createCustomServiceDto := dto.NewCreateCustomService(
				svcName,
				svcType,
				svcCommand,
				svcVersionPtr,
				portBindings,
				autoCreateMapping,
			)

			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

			err := useCase.CreateCustomService(
				servicesQueryRepo,
				servicesCmdRepo,
				mappingQueryRepo,
				mappingCmdRepo,
				vhostQueryRepo,
				createCustomServiceDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "CustomServiceCreated")
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

func (controller *ServicesController) Update() *cobra.Command {
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
			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

			err := useCase.UpdateService(
				servicesQueryRepo,
				servicesCmdRepo,
				mappingQueryRepo,
				mappingCmdRepo,
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

func (controller *ServicesController) Delete() *cobra.Command {
	var nameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteService",
		Run: func(cmd *cobra.Command, args []string) {
			svcName := valueObject.NewServiceNamePanic(nameStr)

			servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
			servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

			err := useCase.DeleteService(
				servicesQueryRepo,
				servicesCmdRepo,
				mappingCmdRepo,
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
