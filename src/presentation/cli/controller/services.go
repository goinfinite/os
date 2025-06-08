package cliController

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	"github.com/spf13/cobra"
)

type ServicesController struct {
	servicesLiaison *liaison.ServicesLiaison
}

func NewServicesController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ServicesController {
	return &ServicesController{
		servicesLiaison: liaison.NewServicesLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (controller *ServicesController) ReadInstalledItems() *cobra.Command {
	var installedItemNameStr, installedItemNatureStr, installedItemTypeStr,
		shouldIncludeMetricsStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadServices",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"shouldIncludeMetrics": shouldIncludeMetricsStr,
			}

			if installedItemNameStr != "" {
				requestBody["name"] = installedItemNameStr
			}

			if installedItemNatureStr != "" {
				requestBody["nature"] = installedItemNatureStr
			}

			if installedItemTypeStr != "" {
				requestBody["type"] = installedItemTypeStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}

			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}

			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}

			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}

			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.servicesLiaison.ReadInstalledItems(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(
		&installedItemNameStr, "installed-item-name", "n", "",
		"InstalledItemName",
	)
	cmd.Flags().StringVarP(
		&installedItemNatureStr, "installed-item-nature", "t", "", "InstalledItemNature",
	)
	cmd.Flags().StringVarP(
		&installedItemTypeStr, "installed-item-type", "u", "", "InstalledItemTypeStr",
	)
	cmd.Flags().StringVarP(
		&shouldIncludeMetricsStr, "should-include-metrics", "s", "false",
		"ShouldIncludeMetrics",
	)
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0,
		"ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "",
		"SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *ServicesController) ReadInstallableItems() *cobra.Command {
	var installedItemNameStr, installedItemNatureStr, installedItemTypeStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get-installables",
		Short: "ReadInstallableServices",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if installedItemNameStr != "" {
				requestBody["name"] = installedItemNameStr
			}

			if installedItemNatureStr != "" {
				requestBody["nature"] = installedItemNatureStr
			}

			if installedItemTypeStr != "" {
				requestBody["type"] = installedItemTypeStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}

			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}

			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}

			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}

			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.servicesLiaison.ReadInstallableItems(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(
		&installedItemNameStr, "installed-item-name", "n", "",
		"InstalledItemName",
	)
	cmd.Flags().StringVarP(
		&installedItemNatureStr, "installed-item-nature", "t", "", "InstalledItemNature",
	)
	cmd.Flags().StringVarP(
		&installedItemTypeStr, "installed-item-type", "u", "", "InstalledItemTypeStr",
	)
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0,
		"ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "",
		"SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *ServicesController) CreateInstallable() *cobra.Command {
	var (
		nameStr, versionStr, startupFileStr, workingDirStr, autoStartBoolStr,
		autoRestartBoolStr, autoCreateMappingBoolStr, mappingHostname, mappingPath,
		mappingUpgradeInsecureRequestsBoolStr string
		envsSlice, portBindingsSlice            []string
		timeoutStartSecsInt, maxStartRetriesInt int
	)

	cmd := &cobra.Command{
		Use:   "create-installable",
		Short: "CreateInstallableService",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name":                           nameStr,
				"autoStart":                      autoStartBoolStr,
				"autoRestart":                    autoRestartBoolStr,
				"autoCreateMapping":              autoCreateMappingBoolStr,
				"mappingUpgradeInsecureRequests": mappingUpgradeInsecureRequestsBoolStr,
			}

			if len(envsSlice) > 0 {
				requestBody["envs"] = sharedHelper.StringSliceValueObjectParser(
					envsSlice, valueObject.NewServiceEnv,
				)
			}
			if versionStr != "" {
				requestBody["version"] = versionStr
			}
			if startupFileStr != "" {
				requestBody["startupFile"] = startupFileStr
			}
			if workingDirStr != "" {
				requestBody["workingDir"] = workingDirStr
			}
			if len(portBindingsSlice) > 0 {
				requestBody["portBindings"] = sharedHelper.StringSliceValueObjectParser(
					portBindingsSlice, valueObject.NewPortBinding,
				)
			}
			if timeoutStartSecsInt != 0 {
				requestBody["timeoutStartSecs"] = uint(timeoutStartSecsInt)
			}
			if maxStartRetriesInt != 0 {
				requestBody["maxStartRetries"] = uint(maxStartRetriesInt)
			}
			if mappingHostname != "" {
				requestBody["mappingHostname"] = mappingHostname
			}
			if mappingPath != "" {
				requestBody["mappingPath"] = mappingPath
			}

			cliHelper.LiaisonResponseWrapper(
				controller.servicesLiaison.CreateInstallable(requestBody, false),
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
	cmd.Flags().StringVarP(&workingDirStr, "working-dir", "w", "", "WorkingDir")
	cmd.Flags().StringSliceVarP(
		&portBindingsSlice, "port-bindings", "p", []string{},
		"PortBindings (port/protocol)",
	)
	cmd.Flags().IntVarP(
		&timeoutStartSecsInt, "timeout-start-secs", "o", 0, "TimeoutStartSecs",
	)
	cmd.Flags().StringVarP(
		&autoStartBoolStr, "auto-start", "s", "true", "AutoStart (true|false)",
	)
	cmd.Flags().StringVarP(
		&autoRestartBoolStr, "auto-restart", "r", "true", "AutoRestart (true|false)",
	)
	cmd.Flags().IntVarP(
		&maxStartRetriesInt, "max-start-retries", "m", 0, "MaxStartRetries",
	)
	cmd.Flags().StringVarP(
		&autoCreateMappingBoolStr, "auto-create-mapping", "a", "true",
		"AutoCreateMapping (true|false)",
	)
	cmd.Flags().StringVarP(
		&mappingHostname, "mapping-hostname", "H", "", "MappingHostname (for AutoCreateMapping)",
	)
	cmd.Flags().StringVarP(
		&mappingPath, "mapping-path", "P", "", "MappingPath (for AutoCreateMapping)",
	)
	cmd.Flags().StringVarP(
		&mappingUpgradeInsecureRequestsBoolStr, "mapping-upgrade-insecure-requests", "u",
		"false", "MappingUpgradeInsecureRequests (true|false)",
	)

	return cmd
}

func (controller *ServicesController) CreateCustom() *cobra.Command {
	var (
		nameStr, typeStr, startCmdStr, versionStr, autoStartBoolStr, autoRestartBoolStr,
		autoCreateMappingBoolStr, mappingHostname, mappingPath,
		mappingUpgradeInsecureRequestsBoolStr string
		envsSlice, portBindingsSlice            []string
		timeoutStartSecsInt, maxStartRetriesInt int
	)

	cmd := &cobra.Command{
		Use:   "create-custom",
		Short: "CreateCustomService",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name":                           nameStr,
				"type":                           typeStr,
				"startCmd":                       startCmdStr,
				"autoStart":                      autoStartBoolStr,
				"autoRestart":                    autoRestartBoolStr,
				"autoCreateMapping":              autoCreateMappingBoolStr,
				"mappingUpgradeInsecureRequests": mappingUpgradeInsecureRequestsBoolStr,
			}

			if len(envsSlice) > 0 {
				requestBody["envs"] = sharedHelper.StringSliceValueObjectParser(
					envsSlice, valueObject.NewServiceEnv,
				)
			}
			if versionStr != "" {
				requestBody["version"] = versionStr
			}
			if len(portBindingsSlice) > 0 {
				requestBody["portBindings"] = sharedHelper.StringSliceValueObjectParser(
					portBindingsSlice, valueObject.NewPortBinding,
				)
			}
			if timeoutStartSecsInt != 0 {
				requestBody["timeoutStartSecs"] = uint(timeoutStartSecsInt)
			}
			if maxStartRetriesInt != 0 {
				requestBody["maxStartRetries"] = uint(maxStartRetriesInt)
			}
			if mappingHostname != "" {
				requestBody["mappingHostname"] = mappingHostname
			}
			if mappingPath != "" {
				requestBody["mappingPath"] = mappingPath
			}

			cliHelper.LiaisonResponseWrapper(
				controller.servicesLiaison.CreateCustom(requestBody),
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
	cmd.Flags().StringVarP(
		&autoStartBoolStr, "auto-start", "s", "true", "AutoStart (true|false)",
	)
	cmd.Flags().IntVarP(
		&maxStartRetriesInt, "max-start-retries", "m", 0, "MaxStartRetries",
	)
	cmd.Flags().StringVarP(
		&autoRestartBoolStr, "auto-restart", "r", "true", "AutoRestart (true|false)",
	)
	cmd.Flags().StringVarP(
		&autoCreateMappingBoolStr, "auto-create-mapping", "a", "true",
		"AutoCreateMapping (true|false)",
	)
	cmd.Flags().StringVarP(
		&mappingHostname, "mapping-hostname", "H", "", "MappingHostname (for AutoCreateMapping)",
	)
	cmd.Flags().StringVarP(
		&mappingPath, "mapping-path", "P", "", "MappingPath (for AutoCreateMapping)",
	)
	cmd.Flags().StringVarP(
		&mappingUpgradeInsecureRequestsBoolStr, "mapping-upgrade-insecure-requests", "u",
		"false", "MappingUpgradeInsecureRequests (true|false) (for AutoCreateMapping)",
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

			cliHelper.LiaisonResponseWrapper(
				controller.servicesLiaison.Update(requestBody),
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

			cliHelper.LiaisonResponseWrapper(
				controller.servicesLiaison.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "ServiceName")
	cmd.MarkFlagRequired("name")
	return cmd
}
