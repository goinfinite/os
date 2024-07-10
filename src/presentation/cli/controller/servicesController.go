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
			servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)
			servicesList, err := useCase.ReadServicesWithMetrics(servicesQueryRepo)
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
			servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)
			servicesList, err := useCase.ReadInstallableServices(servicesQueryRepo)
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
			svcName, err := valueObject.NewServiceName(nameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			var svcVersionPtr *valueObject.ServiceVersion
			if versionStr != "" {
				version, err := valueObject.NewServiceVersion(versionStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				svcVersionPtr = &version
			}

			var startupFilePtr *valueObject.UnixFilePath
			if startupFileStr != "" {
				startupFile, err := valueObject.NewUnixFilePath(startupFileStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				startupFilePtr = &startupFile
			}

			var portBindings []valueObject.PortBinding
			for _, portBinding := range portBindingsSlice {
				svcPortBinding, err := valueObject.NewPortBinding(portBinding)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				portBindings = append(portBindings, svcPortBinding)
			}

			createInstallableServiceDto := dto.NewCreateInstallableService(
				svcName, []valueObject.ServiceEnv{}, portBindings, svcVersionPtr,
				startupFilePtr, nil, nil, nil, nil, &autoCreateMapping,
			)

			servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)
			servicesCmdRepo := servicesInfra.NewServicesCmdRepo(controller.persistentDbSvc)
			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

			err = useCase.CreateInstallableService(
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
		&autoCreateMapping, "auto-create-mapping", "a", true, "AutoCreateMapping",
	)
	return cmd
}

func (controller *ServicesController) CreateCustom() *cobra.Command {
	var nameStr string
	var typeStr string
	var startCmdStr string
	var versionStr string
	var portBindingsSlice []string
	var autoCreateMapping bool

	cmd := &cobra.Command{
		Use:   "create-custom",
		Short: "CreateCustomService",
		Run: func(cmd *cobra.Command, args []string) {
			serviceName, err := valueObject.NewServiceName(nameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			svcType, err := valueObject.NewServiceType(typeStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			startCmd, err := valueObject.NewUnixCommand(startCmdStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			var svcVersionPtr *valueObject.ServiceVersion
			if versionStr != "" {
				version, err := valueObject.NewServiceVersion(versionStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				svcVersionPtr = &version
			}

			var portBindings []valueObject.PortBinding
			for _, portBinding := range portBindingsSlice {
				svcPortBinding, err := valueObject.NewPortBinding(portBinding)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				portBindings = append(portBindings, svcPortBinding)
			}

			createCustomServiceDto := dto.NewCreateCustomService(
				serviceName, svcType, startCmd, []valueObject.ServiceEnv{}, portBindings,
				nil, nil, nil, nil, nil, svcVersionPtr, nil, nil, nil, nil, nil, nil, nil, nil,
				&autoCreateMapping,
			)

			servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)
			servicesCmdRepo := servicesInfra.NewServicesCmdRepo(controller.persistentDbSvc)
			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

			err = useCase.CreateCustomService(
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
	cmd.Flags().StringVarP(&startCmdStr, "start-command", "c", "", "StartCommand")
	cmd.MarkFlagRequired("start-command")
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
	var startCmdStr string
	var statusStr string
	var versionStr string
	var startupFileStr string
	var portBindingsSlice []string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateService",
		Run: func(cmd *cobra.Command, args []string) {
			svcName, err := valueObject.NewServiceName(nameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			var svcTypePtr *valueObject.ServiceType
			if typeStr != "" {
				svcType, err := valueObject.NewServiceType(typeStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				svcTypePtr = &svcType
			}

			var svcStatusPtr *valueObject.ServiceStatus
			if statusStr != "" {
				svcStatus := valueObject.NewServiceStatusPanic(statusStr)
				svcStatusPtr = &svcStatus
			}

			var startCmdPtr *valueObject.UnixCommand
			if startCmdStr != "" {
				startCmd, err := valueObject.NewUnixCommand(startCmdStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				startCmdPtr = &startCmd
			}

			var svcVersionPtr *valueObject.ServiceVersion
			if versionStr != "" {
				version, err := valueObject.NewServiceVersion(versionStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				svcVersionPtr = &version
			}

			var startupFilePtr *valueObject.UnixFilePath
			if startupFileStr != "" {
				startupFile, err := valueObject.NewUnixFilePath(startupFileStr)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				startupFilePtr = &startupFile
			}

			var portBindings []valueObject.PortBinding
			for _, portBinding := range portBindingsSlice {
				svcPortBinding, err := valueObject.NewPortBinding(portBinding)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}
				portBindings = append(portBindings, svcPortBinding)
			}

			updateSvcDto := dto.NewUpdateService(
				svcName, svcTypePtr, svcVersionPtr, svcStatusPtr, startCmdPtr, []valueObject.ServiceEnv{},
				portBindings, nil, nil, nil, nil, nil, startupFilePtr, nil, nil, nil, nil,
				nil, nil, nil, nil,
			)

			servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)
			servicesCmdRepo := servicesInfra.NewServicesCmdRepo(controller.persistentDbSvc)
			mappingQueryRepo := mappingInfra.NewMappingQueryRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

			err = useCase.UpdateService(
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
	cmd.Flags().StringVarP(&startCmdStr, "start-command", "c", "", "StartCommand")
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
			svcName, err := valueObject.NewServiceName(nameStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			servicesQueryRepo := servicesInfra.NewServicesQueryRepo(controller.persistentDbSvc)
			servicesCmdRepo := servicesInfra.NewServicesCmdRepo(controller.persistentDbSvc)
			mappingCmdRepo := mappingInfra.NewMappingCmdRepo(controller.persistentDbSvc)

			err = useCase.DeleteService(
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
