package cliController

import (
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type ServicesController struct {
	persistentDbService *internalDbInfra.PersistentDatabaseService
	serviceServices     *service.ServicesService
}

func NewServicesController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *ServicesController {
	return &ServicesController{
		persistentDbService: persistentDbService,
		serviceServices:     service.NewServicesService(persistentDbService),
	}
}

func (controller *ServicesController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetServices",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.serviceServices.Read())
		},
	}

	return cmd
}

func (controller *ServicesController) ReadInstallables() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-installables",
		Short: "GetInstallableServices",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(
				controller.serviceServices.ReadInstallables(),
			)
		},
	}

	return cmd
}

func (controller *ServicesController) CreateInstallable() *cobra.Command {
	var nameStr, versionStr, startupFileStr string
	var portBindingsSlice []string
	var autoCreateMapping bool

	cmd := &cobra.Command{
		Use:   "create-installable",
		Short: "CreateInstallableService",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name":              nameStr,
				"autoCreateMapping": autoCreateMapping,
			}

			if versionStr != "" {
				requestBody["version"] = versionStr
			}

			if startupFileStr != "" {
				requestBody["startupFile"] = startupFileStr
			}

			if len(portBindingsSlice) > 0 {
				requestBody["portBindings"] = portBindingsSlice
			}

			cliHelper.ServiceResponseWrapper(
				controller.serviceServices.CreateInstallable(requestBody, false),
			)
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
	var nameStr, typeStr, startCmdStr, versionStr string
	var portBindingsSlice []string
	var autoCreateMapping bool

	cmd := &cobra.Command{
		Use:   "create-custom",
		Short: "CreateCustomService",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name":              nameStr,
				"type":              typeStr,
				"startCmd":          startCmdStr,
				"autoCreateMapping": autoCreateMapping,
			}

			if versionStr != "" {
				requestBody["version"] = versionStr
			}

			if len(portBindingsSlice) > 0 {
				requestBody["portBindings"] = portBindingsSlice
			}

			cliHelper.ServiceResponseWrapper(
				controller.serviceServices.CreateCustom(requestBody),
			)
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
	var nameStr, typeStr, startCmdStr, statusStr, versionStr, startupFileStr string
	var portBindingsSlice []string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateService",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name": nameStr,
			}

			if typeStr != "" {
				requestBody["type"] = typeStr
			}

			if startCmdStr != "" {
				requestBody["startCmd"] = startCmdStr
			}

			if statusStr != "" {
				requestBody["status"] = statusStr
			}

			if versionStr != "" {
				requestBody["version"] = versionStr
			}

			if startupFileStr != "" {
				requestBody["startupFile"] = startupFileStr
			}

			if len(portBindingsSlice) > 0 {
				requestBody["portBindings"] = portBindingsSlice
			}

			cliHelper.ServiceResponseWrapper(
				controller.serviceServices.Update(requestBody),
			)
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
			requestBody := map[string]interface{}{
				"name": nameStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.serviceServices.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	return cmd
}
