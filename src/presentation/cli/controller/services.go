package cliController

import (
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type ServicesController struct {
	servicesService *service.ServicesService
}

func NewServicesController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *ServicesController {
	return &ServicesController{
		servicesService: service.NewServicesService(persistentDbService),
	}
}

func (controller *ServicesController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadServices",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.servicesService.Read())
		},
	}

	return cmd
}

func (controller *ServicesController) ReadInstallables() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-installables",
		Short: "ReadInstallableServices",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(
				controller.servicesService.ReadInstallables(),
			)
		},
	}

	return cmd
}

func (controller *ServicesController) CreateInstallable() *cobra.Command {
	var nameStr, versionStr, startupFileStr string
	var envsSlice, portBindingsSlice []string
	var autoStart, autoRestart, autoCreateMapping bool
	var timeoutStartSecsInt, maxStartRetriesInt int

	cmd := &cobra.Command{
		Use:   "create-installable",
		Short: "CreateInstallableService",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name":              nameStr,
				"autoStart":         autoStart,
				"autoRestart":       autoRestart,
				"autoCreateMapping": autoCreateMapping,
			}

			if len(envsSlice) > 0 {
				requestBody["envs"] = envsSlice
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

			if timeoutStartSecsInt != 0 {
				requestBody["timeoutStartSecs"] = uint(timeoutStartSecsInt)
			}

			if maxStartRetriesInt != 0 {
				requestBody["maxStartRetries"] = uint(maxStartRetriesInt)
			}

			cliHelper.ServiceResponseWrapper(
				controller.servicesService.CreateInstallable(requestBody, false),
			)
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringSliceVarP(
		&envsSlice, "envs", "e", []string{}, "Envs (name=value)",
	)
	cmd.Flags().StringVarP(&versionStr, "version", "v", "", "ServiceVersion")
	cmd.Flags().StringVarP(&startupFileStr, "startup-file", "f", "", "StartupFile")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "p", []string{},
		"PortBindings (port/protocol)",
	)
	cmd.Flags().IntVarP(
		&timeoutStartSecsInt, "timeout-start-secs", "o", 0, "TimeoutStartSecs",
	)
	cmd.Flags().BoolVarP(
		&autoStart, "auto-start", "s", true, "AutoStart",
	)
	cmd.Flags().IntVarP(
		&maxStartRetriesInt, "max-start-retries", "m", 0, "MaxStartRetries",
	)
	cmd.Flags().BoolVarP(
		&autoRestart, "auto-restart", "r", true, "AutoRestart",
	)
	cmd.Flags().BoolVarP(
		&autoCreateMapping, "auto-create-mapping", "a", true, "AutoCreateMapping",
	)
	return cmd
}

func (controller *ServicesController) CreateCustom() *cobra.Command {
	var nameStr, typeStr, startCmdStr, versionStr string
	var envsSlice, portBindingsSlice []string
	var autoStart, autoRestart, autoCreateMapping bool
	var timeoutStartSecsInt, maxStartRetriesInt int

	cmd := &cobra.Command{
		Use:   "create-custom",
		Short: "CreateCustomService",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name":              nameStr,
				"type":              typeStr,
				"startCmd":          startCmdStr,
				"autoStart":         autoStart,
				"autoRestart":       autoRestart,
				"autoCreateMapping": autoCreateMapping,
			}

			if len(envsSlice) > 0 {
				requestBody["envs"] = envsSlice
			}

			if versionStr != "" {
				requestBody["version"] = versionStr
			}

			if len(portBindingsSlice) > 0 {
				requestBody["portBindings"] = portBindingsSlice
			}

			if timeoutStartSecsInt != 0 {
				requestBody["timeoutStartSecs"] = uint(timeoutStartSecsInt)
			}

			if maxStartRetriesInt != 0 {
				requestBody["maxStartRetries"] = uint(maxStartRetriesInt)
			}

			cliHelper.ServiceResponseWrapper(
				controller.servicesService.CreateCustom(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(
		&typeStr, "type", "t", "", "ServiceType (application|database|runtime|other)",
	)
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&startCmdStr, "start-command", "c", "", "StartCommand")
	cmd.MarkFlagRequired("start-command")
	cmd.Flags().StringSliceVarP(
		&envsSlice, "envs", "e", []string{}, "Envs (name=value)",
	)
	cmd.Flags().StringVarP(&versionStr, "version", "v", "", "ServiceVersion")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "p", []string{},
		"PortBindings (port/protocol)",
	)
	cmd.Flags().IntVarP(
		&timeoutStartSecsInt, "timeout-start-secs", "o", 0, "TimeoutStartSecs",
	)
	cmd.Flags().BoolVarP(
		&autoStart, "auto-start", "s", true, "AutoStart",
	)
	cmd.Flags().IntVarP(
		&maxStartRetriesInt, "max-start-retries", "m", 0, "MaxStartRetries",
	)
	cmd.Flags().BoolVarP(
		&autoRestart, "auto-restart", "r", true, "AutoRestart",
	)
	cmd.Flags().BoolVarP(
		&autoCreateMapping, "auto-create-mapping", "a", true, "AutoCreateMapping",
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
				controller.servicesService.Update(requestBody),
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
		&portBindingsSlice, "port-bindings", "p", []string{},
		"PortBindings (port/protocol)",
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
				controller.servicesService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	return cmd
}
