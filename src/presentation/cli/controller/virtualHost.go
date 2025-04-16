package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type VirtualHostController struct {
	virtualHostService *service.VirtualHostService
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		virtualHostService: service.NewVirtualHostService(persistentDbSvc, trailDbSvc),
	}
}

func (controller *VirtualHostController) Read() *cobra.Command {
	var hostnameStr, typeStr, rootDirectoryStr, parentHostnameStr, withMappingsBoolStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadVirtualHosts",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if hostnameStr != "" {
				requestBody["hostname"] = hostnameStr
			}
			if typeStr != "" {
				requestBody["type"] = typeStr
			}
			if rootDirectoryStr != "" {
				requestBody["rootDirectory"] = rootDirectoryStr
			}
			if parentHostnameStr != "" {
				requestBody["parentHostname"] = parentHostnameStr
			}
			if withMappingsBoolStr != "" {
				requestBody["withMappings"] = withMappingsBoolStr
			}

			requestBody = cliHelper.PaginationParser(
				requestBody, paginationPageNumberUint32, paginationItemsPerPageUint16,
				paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr,
			)

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.Read(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.Flags().StringVarP(&typeStr, "type", "t", "", "VirtualHostType")
	cmd.Flags().StringVarP(&rootDirectoryStr, "root", "r", "", "RootDirectory")
	cmd.Flags().StringVarP(&parentHostnameStr, "parent", "p", "", "ParentHostname")
	cmd.Flags().StringVarP(
		&withMappingsBoolStr, "with-mappings", "w", "false", "WithMappings (true|false)",
	)
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "o", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "j", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "x", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)
	return cmd
}

func (controller *VirtualHostController) Create() *cobra.Command {
	var hostnameStr, typeStr, parentHostnameStr, isWildcardBoolStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"hostname":   hostnameStr,
				"isWildcard": isWildcardBoolStr,
			}

			if typeStr != "" {
				requestBody["type"] = typeStr
			}

			if parentHostnameStr != "" {
				requestBody["parentHostname"] = parentHostnameStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(
		&typeStr, "type", "t", "", "VirtualHostType (top-level|subdomain|alias)",
	)
	cmd.Flags().StringVarP(&parentHostnameStr, "parent", "p", "", "ParentHostname")
	cmd.Flags().StringVarP(
		&isWildcardBoolStr, "is-wildcard", "w", "false", "IsWildcard (true|false)",
	)
	return cmd
}

func (controller *VirtualHostController) Update() *cobra.Command {
	var hostnameStr, isWildcardBoolStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"hostname":   hostnameStr,
				"isWildcard": isWildcardBoolStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.Update(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(
		&isWildcardBoolStr, "is-wildcard", "w", "false", "IsWildcard (true|false)",
	)
	return cmd
}

func (controller *VirtualHostController) Delete() *cobra.Command {
	var hostnameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteVirtualHost",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"hostname": hostnameStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.MarkFlagRequired("hostname")
	return cmd
}

func (controller *VirtualHostController) ReadWithMappings() *cobra.Command {
	var hostnameStr, typeStr, rootDirectoryStr, parentHostnameStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadVirtualHostsWithMappings",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"withMappings": true,
			}

			if hostnameStr != "" {
				requestBody["hostname"] = hostnameStr
			}
			if typeStr != "" {
				requestBody["type"] = typeStr
			}
			if rootDirectoryStr != "" {
				requestBody["rootDirectory"] = rootDirectoryStr
			}
			if parentHostnameStr != "" {
				requestBody["parentHostname"] = parentHostnameStr
			}

			requestBody = cliHelper.PaginationParser(
				requestBody, paginationPageNumberUint32, paginationItemsPerPageUint16,
				paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr,
			)

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.Read(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostHostname")
	cmd.Flags().StringVarP(&typeStr, "type", "t", "", "VirtualHostType")
	cmd.Flags().StringVarP(&rootDirectoryStr, "root", "r", "", "RootDirectory")
	cmd.Flags().StringVarP(&parentHostnameStr, "parent", "p", "", "ParentHostname")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "o", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "j", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "x", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)
	return cmd
}

func (controller *VirtualHostController) CreateMapping() *cobra.Command {
	var (
		hostnameStr, pathStr, matchPatternStr, targetTypeStr, targetValueStr,
		shouldUpgradeInsecureRequestsBoolStr string
		targetHttpResponseCodeUint, mappingSecurityRuleIdUint uint
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateVirtualHostMapping",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"hostname":                      hostnameStr,
				"path":                          pathStr,
				"targetType":                    targetTypeStr,
				"shouldUpgradeInsecureRequests": shouldUpgradeInsecureRequestsBoolStr,
			}

			if matchPatternStr != "" {
				requestBody["matchPattern"] = matchPatternStr
			}

			if targetValueStr != "" {
				requestBody["targetValue"] = targetValueStr
			}

			if targetHttpResponseCodeUint != 0 {
				requestBody["targetHttpResponseCode"] = targetHttpResponseCodeUint
			}

			if mappingSecurityRuleIdUint != 0 {
				requestBody["mappingSecurityRuleId"] = mappingSecurityRuleIdUint
			}

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.CreateMapping(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "Hostname")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().StringVarP(&pathStr, "path", "p", "", "MappingPath")
	cmd.MarkFlagRequired("path")
	cmd.Flags().StringVarP(
		&matchPatternStr, "match", "m", "",
		"MatchPattern (begins-with|contains|ends-with)",
	)
	cmd.Flags().StringVarP(
		&targetTypeStr, "type", "t", "",
		"MappingTargetType (url|service|response-code|inline-html|static-files)",
	)
	cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&targetValueStr, "value", "v", "", "MappingTargetValue")
	cmd.Flags().UintVarP(
		&targetHttpResponseCodeUint, "response-code", "r", 0, "TargetHttpResponseCode",
	)
	cmd.Flags().StringVarP(
		&shouldUpgradeInsecureRequestsBoolStr, "should-upgrade-insecure-requests", "u",
		"false", "ShouldUpgradeInsecureRequests (true|false)",
	)
	cmd.Flags().UintVarP(
		&mappingSecurityRuleIdUint, "mapping-security-rule-id", "s", 0, "MappingSecurityRuleId",
	)
	return cmd
}

func (controller *VirtualHostController) DeleteMapping() *cobra.Command {
	var mappingIdUint uint

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteVirtualHostMapping",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": mappingIdUint,
			}

			cliHelper.ServiceResponseWrapper(
				controller.virtualHostService.DeleteMapping(requestBody),
			)
		},
	}

	cmd.Flags().UintVarP(&mappingIdUint, "id", "i", 0, "MappingId")
	cmd.MarkFlagRequired("id")
	return cmd
}
