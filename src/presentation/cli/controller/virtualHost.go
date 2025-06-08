package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/spf13/cobra"
)

type VirtualHostController struct {
	virtualHostLiaison *liaison.VirtualHostLiaison
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		virtualHostLiaison: liaison.NewVirtualHostLiaison(persistentDbSvc, trailDbSvc),
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

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.Read(requestBody),
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

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.Create(requestBody),
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

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.Update(requestBody),
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

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.Delete(requestBody),
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

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.Read(requestBody),
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

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.CreateMapping(requestBody),
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

func (controller *VirtualHostController) UpdateMapping() *cobra.Command {
	var (
		mappingIdUint uint
		pathStr, matchPatternStr, targetTypeStr, targetValueStr,
		shouldUpgradeInsecureRequestsBoolStr string
		targetHttpResponseCodeUint, mappingSecurityRuleIdUint uint
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateVirtualHostMapping",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": mappingIdUint,
			}

			if pathStr != "" {
				requestBody["path"] = pathStr
			}

			if matchPatternStr != "" {
				requestBody["matchPattern"] = matchPatternStr
			}

			if targetTypeStr != "" {
				requestBody["targetType"] = targetTypeStr
			}

			if targetValueStr != "" {
				requestBody["targetValue"] = targetValueStr
			}

			if targetHttpResponseCodeUint != 0 {
				requestBody["targetHttpResponseCode"] = targetHttpResponseCodeUint
			}

			if shouldUpgradeInsecureRequestsBoolStr != "" {
				requestBody["shouldUpgradeInsecureRequests"] = shouldUpgradeInsecureRequestsBoolStr
			}

			if mappingSecurityRuleIdUint != 0 {
				requestBody["mappingSecurityRuleId"] = mappingSecurityRuleIdUint
			}

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.UpdateMapping(requestBody),
			)
		},
	}

	cmd.Flags().UintVarP(&mappingIdUint, "id", "i", 0, "MappingId")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&pathStr, "path", "p", "", "MappingPath")
	cmd.Flags().StringVarP(
		&matchPatternStr, "match", "m", "",
		"MatchPattern (begins-with|contains|ends-with)",
	)
	cmd.Flags().StringVarP(
		&targetTypeStr, "type", "t", "",
		"MappingTargetType (url|service|response-code|inline-html|static-files)",
	)
	cmd.Flags().StringVarP(&targetValueStr, "value", "v", "", "MappingTargetValue")
	cmd.Flags().UintVarP(
		&targetHttpResponseCodeUint, "response-code", "r", 0, "TargetHttpResponseCode",
	)
	cmd.Flags().StringVarP(
		&shouldUpgradeInsecureRequestsBoolStr, "should-upgrade-insecure-requests", "u",
		"", "ShouldUpgradeInsecureRequests (true|false)",
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

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.DeleteMapping(requestBody),
			)
		},
	}

	cmd.Flags().UintVarP(&mappingIdUint, "id", "i", 0, "MappingId")
	cmd.MarkFlagRequired("id")
	return cmd
}

func (controller *VirtualHostController) ReadMappingSecurityRules() *cobra.Command {
	var ruleIdUint uint
	var ruleNameStr, allowedIpStr, blockedIpStr string
	var createdBeforeAtInt, createdAfterAtInt int64
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadMappingSecurityRules",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if ruleIdUint != 0 {
				requestBody["id"] = ruleIdUint
			}
			if ruleNameStr != "" {
				requestBody["name"] = ruleNameStr
			}
			if allowedIpStr != "" {
				requestBody["allowedIp"] = allowedIpStr
			}
			if blockedIpStr != "" {
				requestBody["blockedIp"] = blockedIpStr
			}
			if createdBeforeAtInt != 0 {
				requestBody["createdBeforeAt"] = createdBeforeAtInt
			}
			if createdAfterAtInt != 0 {
				requestBody["createdAfterAt"] = createdAfterAtInt
			}

			requestBody = cliHelper.PaginationParser(
				requestBody, paginationPageNumberUint32, paginationItemsPerPageUint16,
				paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr,
			)

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.ReadMappingSecurityRules(requestBody),
			)
		},
	}

	cmd.Flags().UintVarP(&ruleIdUint, "id", "i", 0, "MappingSecurityRuleId")
	cmd.Flags().StringVarP(&ruleNameStr, "name", "n", "", "MappingSecurityRuleName")
	cmd.Flags().StringVarP(&allowedIpStr, "allowed-ip", "a", "", "AllowedIpAddress")
	cmd.Flags().StringVarP(&blockedIpStr, "blocked-ip", "b", "", "BlockedIpAddress")
	cmd.Flags().Int64VarP(
		&createdBeforeAtInt, "created-before", "e", 0, "CreatedBeforeAt (Unix timestamp)",
	)
	cmd.Flags().Int64VarP(
		&createdAfterAtInt, "created-after", "f", 0, "CreatedAfterAt (Unix timestamp)",
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

func (controller *VirtualHostController) CreateMappingSecurityRule() *cobra.Command {
	var nameStr, descriptionStr string
	var allowedIpsSlice, blockedIpsSlice []string
	var rpsSoftLimitPerIpUint, rpsHardLimitPerIpUint, responseCodeOnMaxRequestsUint uint
	var maxConnectionsPerIpUint uint
	var bandwidthBpsLimitPerConnectionUint64, bandwidthLimitOnlyAfterBytesUint64 uint64
	var responseCodeOnMaxConnectionsUint uint

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateMappingSecurityRule",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"name": nameStr,
			}

			if descriptionStr != "" {
				requestBody["description"] = descriptionStr
			}

			if len(allowedIpsSlice) > 0 {
				requestBody["allowedIps"] = tkPresentation.StringSliceValueObjectParser(
					allowedIpsSlice, tkValueObject.NewCidrBlock,
				)
			}

			if len(blockedIpsSlice) > 0 {
				requestBody["blockedIps"] = tkPresentation.StringSliceValueObjectParser(
					blockedIpsSlice, tkValueObject.NewCidrBlock,
				)
			}

			if rpsSoftLimitPerIpUint != 0 {
				requestBody["rpsSoftLimitPerIp"] = rpsSoftLimitPerIpUint
			}

			if rpsHardLimitPerIpUint != 0 {
				requestBody["rpsHardLimitPerIp"] = rpsHardLimitPerIpUint
			}

			if responseCodeOnMaxRequestsUint != 0 {
				requestBody["responseCodeOnMaxRequests"] = responseCodeOnMaxRequestsUint
			}

			if maxConnectionsPerIpUint != 0 {
				requestBody["maxConnectionsPerIp"] = maxConnectionsPerIpUint
			}

			if bandwidthBpsLimitPerConnectionUint64 != 0 {
				requestBody["bandwidthBpsLimitPerConnection"] = bandwidthBpsLimitPerConnectionUint64
			}

			if bandwidthLimitOnlyAfterBytesUint64 != 0 {
				requestBody["bandwidthLimitOnlyAfterBytes"] = bandwidthLimitOnlyAfterBytesUint64
			}

			if responseCodeOnMaxConnectionsUint != 0 {
				requestBody["responseCodeOnMaxConnections"] = responseCodeOnMaxConnectionsUint
			}

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.CreateMappingSecurityRule(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "MappingSecurityRuleName")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(
		&descriptionStr, "description", "d", "", "MappingSecurityRuleDescription",
	)
	cmd.Flags().StringSliceVarP(
		&allowedIpsSlice, "allowed-ips", "a", []string{},
		"AllowedIps (comma-separated CIDR blocks)",
	)
	cmd.Flags().StringSliceVarP(
		&blockedIpsSlice, "blocked-ips", "b", []string{},
		"BlockedIps (comma-separated CIDR blocks)",
	)
	cmd.Flags().UintVarP(
		&rpsSoftLimitPerIpUint, "rps-soft-limit", "S", 0, "RpsSoftLimitPerIp",
	)
	cmd.Flags().UintVarP(
		&rpsHardLimitPerIpUint, "rps-hard-limit", "H", 0, "RpsHardLimitPerIp",
	)
	cmd.Flags().UintVarP(
		&responseCodeOnMaxRequestsUint, "response-code-requests", "r", 0,
		"ResponseCodeOnMaxRequests",
	)
	cmd.Flags().UintVarP(
		&maxConnectionsPerIpUint, "max-connections", "m", 0, "MaxConnectionsPerIp",
	)
	cmd.Flags().Uint64VarP(
		&bandwidthBpsLimitPerConnectionUint64, "bandwidth-limit", "l", 0,
		"BandwidthBpsLimitPerConnection",
	)
	cmd.Flags().Uint64VarP(
		&bandwidthLimitOnlyAfterBytesUint64, "bandwidth-limit-after", "f", 0,
		"BandwidthLimitOnlyAfterBytes",
	)
	cmd.Flags().UintVarP(
		&responseCodeOnMaxConnectionsUint, "response-code-connections", "c", 0,
		"ResponseCodeOnMaxConnections",
	)
	return cmd
}

func (controller *VirtualHostController) UpdateMappingSecurityRule() *cobra.Command {
	var ruleIdUint uint
	var nameStr, descriptionStr string
	var allowedIpsSlice, blockedIpsSlice []string
	var rpsSoftLimitPerIpUint, rpsHardLimitPerIpUint, responseCodeOnMaxRequestsUint uint
	var maxConnectionsPerIpUint uint
	var bandwidthBpsLimitPerConnectionUint64, bandwidthLimitOnlyAfterBytesUint64 uint64
	var responseCodeOnMaxConnectionsUint uint

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateMappingSecurityRule",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": ruleIdUint,
			}

			if nameStr != "" {
				requestBody["name"] = nameStr
			}

			if descriptionStr != "" {
				requestBody["description"] = descriptionStr
			}

			if len(allowedIpsSlice) > 0 {
				requestBody["allowedIps"] = tkPresentation.StringSliceValueObjectParser(
					allowedIpsSlice, tkValueObject.NewCidrBlock,
				)
			}

			if len(blockedIpsSlice) > 0 {
				requestBody["blockedIps"] = tkPresentation.StringSliceValueObjectParser(
					blockedIpsSlice, tkValueObject.NewCidrBlock,
				)
			}

			if rpsSoftLimitPerIpUint != 0 {
				requestBody["rpsSoftLimitPerIp"] = rpsSoftLimitPerIpUint
			}

			if rpsHardLimitPerIpUint != 0 {
				requestBody["rpsHardLimitPerIp"] = rpsHardLimitPerIpUint
			}

			if responseCodeOnMaxRequestsUint != 0 {
				requestBody["responseCodeOnMaxRequests"] = responseCodeOnMaxRequestsUint
			}

			if maxConnectionsPerIpUint != 0 {
				requestBody["maxConnectionsPerIp"] = maxConnectionsPerIpUint
			}

			if bandwidthBpsLimitPerConnectionUint64 != 0 {
				requestBody["bandwidthBpsLimitPerConnection"] = bandwidthBpsLimitPerConnectionUint64
			}

			if bandwidthLimitOnlyAfterBytesUint64 != 0 {
				requestBody["bandwidthLimitOnlyAfterBytes"] = bandwidthLimitOnlyAfterBytesUint64
			}

			if responseCodeOnMaxConnectionsUint != 0 {
				requestBody["responseCodeOnMaxConnections"] = responseCodeOnMaxConnectionsUint
			}

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.UpdateMappingSecurityRule(requestBody),
			)
		},
	}

	cmd.Flags().UintVarP(&ruleIdUint, "id", "i", 0, "MappingSecurityRuleId")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&nameStr, "name", "n", "", "MappingSecurityRuleName")
	cmd.Flags().StringVarP(
		&descriptionStr, "description", "d", "", "MappingSecurityRuleDescription",
	)
	cmd.Flags().StringSliceVarP(
		&allowedIpsSlice, "allowed-ips", "a", []string{},
		"AllowedIps (comma-separated CIDR blocks)",
	)
	cmd.Flags().StringSliceVarP(
		&blockedIpsSlice, "blocked-ips", "b", []string{},
		"BlockedIps (comma-separated CIDR blocks)",
	)
	cmd.Flags().UintVarP(
		&rpsSoftLimitPerIpUint, "rps-soft-limit", "S", 0, "RpsSoftLimitPerIp",
	)
	cmd.Flags().UintVarP(
		&rpsHardLimitPerIpUint, "rps-hard-limit", "H", 0, "RpsHardLimitPerIp",
	)
	cmd.Flags().UintVarP(
		&responseCodeOnMaxRequestsUint, "response-code-requests", "r", 0,
		"ResponseCodeOnMaxRequests",
	)
	cmd.Flags().UintVarP(
		&maxConnectionsPerIpUint, "max-connections", "m", 0, "MaxConnectionsPerIp",
	)
	cmd.Flags().Uint64VarP(
		&bandwidthBpsLimitPerConnectionUint64, "bandwidth-limit", "l", 0,
		"BandwidthBpsLimitPerConnection",
	)
	cmd.Flags().Uint64VarP(
		&bandwidthLimitOnlyAfterBytesUint64, "bandwidth-limit-after", "f", 0,
		"BandwidthLimitOnlyAfterBytes",
	)
	cmd.Flags().UintVarP(
		&responseCodeOnMaxConnectionsUint, "response-code-connections", "c", 0,
		"ResponseCodeOnMaxConnections",
	)
	return cmd
}

func (controller *VirtualHostController) DeleteMappingSecurityRule() *cobra.Command {
	var ruleIdUint uint

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteMappingSecurityRule",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": ruleIdUint,
			}

			cliHelper.LiaisonResponseWrapper(
				controller.virtualHostLiaison.DeleteMappingSecurityRule(requestBody),
			)
		},
	}

	cmd.Flags().UintVarP(&ruleIdUint, "id", "i", 0, "MappingSecurityRuleId")
	cmd.MarkFlagRequired("id")
	return cmd
}
