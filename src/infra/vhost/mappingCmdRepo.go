package vhostInfra

import (
	"errors"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
)

type MappingCmdRepo struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	mappingQueryRepo *MappingQueryRepo
}

func NewMappingCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MappingCmdRepo {
	mappingQueryRepo := NewMappingQueryRepo(persistentDbSvc)

	return &MappingCmdRepo{
		persistentDbSvc:  persistentDbSvc,
		mappingQueryRepo: mappingQueryRepo,
	}
}

func (repo *MappingCmdRepo) serviceMappingConfigFactory(
	svcNameStr string,
) (svcMappingConfig string, err error) {
	svcMappingConfig = ""

	serviceName, err := valueObject.NewServiceName(svcNameStr)
	if err != nil {
		return svcMappingConfig, errors.New(err.Error() + ": " + svcNameStr)
	}

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(repo.persistentDbSvc)
	readFirstInstalledServiceRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &serviceName,
	}
	installedService, err := servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledServiceRequestDto,
	)
	if err != nil {
		return svcMappingConfig, err
	}

	protocolPortsMap := map[string]string{}
	for _, svcPortBinding := range installedService.PortBindings {
		svcPortBindingProtocolStr := svcPortBinding.Protocol.String()
		protocolPortsMap[svcPortBindingProtocolStr] = svcPortBinding.Port.String()
	}

	isHttpSupported := protocolPortsMap["http"] != ""
	if isHttpSupported {
		svcMappingConfig += `
	set $protocol "http";
	set $backend "127.0.0.1:` + protocolPortsMap["http"] + `";
`
	}

	isHttpsSupported := protocolPortsMap["https"] != ""
	if isHttpsSupported {
		svcMappingConfig += `
	set $protocol "https";
	set $backend "127.0.0.1:` + protocolPortsMap["https"] + `";
`
	}

	if isHttpSupported && isHttpsSupported {
		svcMappingConfig = `
	set $protocol "http";
	set $backend "127.0.0.1:` + protocolPortsMap["http"] + `";

	if ($scheme = https) {
		set $protocol "https";
		set $backend "127.0.0.1:` + protocolPortsMap["https"] + `";
	}
`
	}

	isHttpOrHttpsSupported := isHttpSupported || isHttpsSupported

	isWsSupported := protocolPortsMap["ws"] != ""
	isWssSupported := protocolPortsMap["wss"] != ""
	if isWsSupported && !isHttpOrHttpsSupported {
		svcMappingConfig += `
	set $protocol "http";
	set $backend "127.0.0.1:` + protocolPortsMap["ws"] + `";
`
	}

	if isWsSupported && !isWssSupported && !isHttpSupported {
		svcMappingConfig += `
	if ($scheme = http) {
		set $protocol "http";
		set $backend "127.0.0.1:` + protocolPortsMap["ws"] + `";
	}
`
	}

	if !isWsSupported && isWssSupported && !isHttpOrHttpsSupported {
		svcMappingConfig += `
	set $protocol "https";
	set $backend "127.0.0.1:` + protocolPortsMap["wss"] + `";
`
	}

	if !isWsSupported && isWssSupported && !isHttpsSupported {
		svcMappingConfig += `
	if ($scheme = https) {
		set $protocol "https";
		set $backend "127.0.0.1:` + protocolPortsMap["wss"] + `";
	}
`
	}

	isWsAndWssSupported := isWsSupported && isWssSupported
	if isWsAndWssSupported && !isHttpOrHttpsSupported {
		svcMappingConfig = `
	set $protocol "http";
	set $backend "127.0.0.1:` + protocolPortsMap["ws"] + `";

	if ($scheme = https) {
		set $protocol "https";
		set $backend "127.0.0.1:` + protocolPortsMap["wss"] + `";
	}
`
	}

	isWsOrWssSupported := isWsSupported || isWssSupported
	if isWsOrWssSupported {
		svcMappingConfig += `
	proxy_http_version 1.1;
	proxy_set_header Upgrade $http_upgrade;
	proxy_set_header Connection "Upgrade";
`
	}

	isHttpOrHttpsSupported = isHttpOrHttpsSupported || isWsOrWssSupported

	isGrpcSupported := protocolPortsMap["grpc"] != ""
	if isGrpcSupported && !isHttpOrHttpsSupported {
		svcMappingConfig += `
	set $protocol "grpc";
	set $backend "127.0.0.1:` + protocolPortsMap["grpc"] + `";
`
	}

	if isGrpcSupported && isHttpOrHttpsSupported {
		svcMappingConfig += `
	if ($scheme = grpc) {
		set $protocol "grpc";
		set $backend "127.0.0.1:` + protocolPortsMap["grpc"] + `";
	}
`
	}

	isGrpcsSupported := protocolPortsMap["grpcs"] != ""
	if isGrpcsSupported && !isHttpOrHttpsSupported {
		svcMappingConfig += `
	set $protocol "grpcs";
	set $backend "127.0.0.1:` + protocolPortsMap["grpcs"] + `";
`
	}

	if isGrpcsSupported && isHttpOrHttpsSupported {
		svcMappingConfig += `
	if ($scheme = grpcs) {
		set $protocol "grpcs";
		set $backend "127.0.0.1:` + protocolPortsMap["grpcs"] + `";
	}
		`
	}

	if isGrpcSupported && !isGrpcsSupported && isHttpOrHttpsSupported {
		svcMappingConfig += `
	grpc_set_header Host $host;
	if ($protocol = grpc) {
		grpc_pass $protocol://$backend;
	}
`
	}

	if !isGrpcSupported && isGrpcsSupported && isHttpOrHttpsSupported {
		svcMappingConfig += `
	grpc_set_header Host $host;
	if ($protocol = grpcs) {
		grpc_pass $protocol://$backend;
	}
`
	}

	isGrpcAndGrpcsSupported := isGrpcSupported && isGrpcsSupported
	if isGrpcAndGrpcsSupported && !isHttpOrHttpsSupported {
		svcMappingConfig = `
	set $protocol "grpc";
	set $backend "127.0.0.1:` + protocolPortsMap["grpc"] + `";

	if ($scheme = grpcs) {
		set $protocol "grpcs";
		set $backend "127.0.0.1:` + protocolPortsMap["grpcs"] + `";
	}
`
	}

	if isGrpcAndGrpcsSupported && isHttpOrHttpsSupported {
		svcMappingConfig += `
	grpc_set_header Host $host;
	if ($protocol = grpc) {
		grpc_pass $protocol://$backend;
	}

	if ($protocol = grpcs) {
		grpc_pass $protocol://$backend;
	}
`
	}

	isGrpcOrGrpcsSupported := isGrpcSupported || isGrpcsSupported
	if isGrpcOrGrpcsSupported && !isHttpOrHttpsSupported {
		svcMappingConfig += `
	grpc_set_header Host $host;
	grpc_pass $protocol://$backend;
`
	}

	if isHttpOrHttpsSupported {
		svcMappingConfig += `
	proxy_pass $protocol://$backend;
	proxy_set_header Host $host;
`
	}

	isTcpSupported := protocolPortsMap["tcp"] != ""
	if isTcpSupported {
		svcMappingConfig += `
	set $protocol "tcp";
	set $backend "127.0.0.1:` + protocolPortsMap["tcp"] + `";
	proxy_pass $protocol://$backend;
`
	}

	svcMappingConfig = strings.Trim(svcMappingConfig, "\n")
	return svcMappingConfig, nil
}

func (repo *MappingCmdRepo) locationUriConfigFactory(
	matchPattern valueObject.MappingMatchPattern,
	path valueObject.MappingPath,
) (locationUri string) {
	matchPatternStr := matchPattern.String()

	modifier := ""
	switch matchPatternStr {
	case "contains", "ends-with":
		modifier = "~"
	case "equals":
		modifier = "="
	}

	pathStr := path.String()
	if matchPatternStr == "ends-with" {
		pathStr += "$"
	}

	locationUri = pathStr
	if modifier != "" {
		locationUri = modifier + " " + pathStr
	}

	return locationUri
}

func (repo *MappingCmdRepo) RecreateMappingFile(
	vhostHostname valueObject.Fqdn,
) error {
	mappingsReadResponse, err := repo.mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Pagination: dto.PaginationUnpaginated,
		Hostname:   &vhostHostname,
	})
	if err != nil {
		return err
	}

	mappingConfigTemplate := `{{- range . -}}
location {{ locationUriConfigFactory .MatchPattern .Path }} {
	{{- if boolPtrDeref .ShouldUpgradeInsecureRequests }}
	if ($scheme = http) {
		return 301 https://$host$request_uri;
	}
	{{- end }}
	{{- if .MappingSecurityRuleId }}
	include ` + infraEnvs.MappingsSecurityRulesConfDir + `/{{ .MappingSecurityRuleId.String }}.embeddable.conf;
	{{- end }}
	{{- if eq .TargetType "response-code" }}
	return {{ .TargetHttpResponseCode }};
	{{- end }}
	{{- if eq .TargetType "url" }}
	return {{ .TargetHttpResponseCode }} {{ .TargetValue }};
	{{- end }}
	{{- if eq .TargetType "inline-html" }}
	add_header Content-Type text/html;
	return {{ .TargetHttpResponseCode }} "{{ .TargetValue }}";
	{{- end }}
	{{- if eq .TargetType "service" }}
{{ serviceMappingConfigFactory .TargetValue.String }}
	{{- end }}
	{{- if eq .TargetType "static-files" }}
	try_files $uri $uri/ index.html?$query_string;
	{{- end }}
}
{{ end }}`

	mappingTemplatePtr := template.New("mappingFile").Funcs(
		template.FuncMap{
			"locationUriConfigFactory":    repo.locationUriConfigFactory,
			"serviceMappingConfigFactory": repo.serviceMappingConfigFactory,
			"boolPtrDeref": func(boolPtr *bool) bool {
				if boolPtr == nil {
					return false
				}
				return *boolPtr
			},
		},
	)

	mappingTemplatePtr, err = mappingTemplatePtr.Parse(mappingConfigTemplate)
	if err != nil {
		return errors.New("TemplateParsingError: " + err.Error())
	}

	var mappingFileContent strings.Builder
	err = mappingTemplatePtr.Execute(&mappingFileContent, mappingsReadResponse.Mappings)
	if err != nil {
		return errors.New("TemplateExecutionError: " + err.Error())
	}

	vhostQueryRepo := NewVirtualHostQueryRepo(repo.persistentDbSvc)
	mappingFilePath, err := vhostQueryRepo.ReadVirtualHostMappingsFilePath(
		vhostHostname,
	)
	if err != nil {
		return errors.New("ReadVirtualHostMappingsFilePathError: " + err.Error())
	}

	shouldOverwrite := true
	return infraHelper.UpdateFile(
		mappingFilePath.String(), mappingFileContent.String(), shouldOverwrite,
	)
}

func (repo *MappingCmdRepo) Create(
	createDto dto.CreateMapping,
) (mappingId valueObject.MappingId, err error) {
	err = infraHelper.ValidateWebServerConfig()
	if err != nil {
		return mappingId, err
	}

	isServiceMapping := createDto.TargetType.String() == "service"
	isPhpServiceMapping := isServiceMapping && createDto.TargetValue.String() == "php-webserver"
	if isPhpServiceMapping {
		runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(repo.persistentDbSvc)
		err := runtimeCmdRepo.CreatePhpVirtualHost(createDto.Hostname)
		if err != nil {
			return mappingId, err
		}
	}

	var targetValuePtr *string
	if createDto.TargetValue != nil {
		targetValueStr := createDto.TargetValue.String()
		targetValuePtr = &targetValueStr
	}

	var targetHttpResponseCodePtr *string
	if createDto.TargetHttpResponseCode != nil {
		targetHttpResponseCodeStr := createDto.TargetHttpResponseCode.String()
		targetHttpResponseCodePtr = &targetHttpResponseCodeStr
	}

	var mappingSecurityRuleIdPtr *uint64
	if createDto.MappingSecurityRuleId != nil {
		ruleIdUint64 := createDto.MappingSecurityRuleId.Uint64()
		mappingSecurityRuleIdPtr = &ruleIdUint64
	}

	mappingModel := dbModel.Mapping{
		Hostname:                      createDto.Hostname.String(),
		Path:                          createDto.Path.String(),
		MatchPattern:                  createDto.MatchPattern.String(),
		TargetType:                    createDto.TargetType.String(),
		TargetValue:                   targetValuePtr,
		TargetHttpResponseCode:        targetHttpResponseCodePtr,
		ShouldUpgradeInsecureRequests: createDto.ShouldUpgradeInsecureRequests,
		MappingSecurityRuleID:         mappingSecurityRuleIdPtr,
	}
	err = repo.persistentDbSvc.Handler.Create(&mappingModel).Error
	if err != nil {
		return mappingId, err
	}

	mappingId, err = valueObject.NewMappingId(mappingModel.ID)
	if err != nil {
		return mappingId, err
	}

	err = repo.RecreateMappingFile(createDto.Hostname)
	if err != nil {
		return mappingId, errors.New("RecreateMappingFileError: " + err.Error())
	}

	err = infraHelper.ValidateWebServerConfig()
	if err != nil {
		err = repo.persistentDbSvc.Handler.Delete(&mappingModel).Error
		if err != nil {
			return mappingId, err
		}

		err = repo.RecreateMappingFile(createDto.Hostname)
		if err != nil {
			return mappingId, errors.New("RecreateMappingFileError: " + err.Error())
		}
	}

	return mappingId, infraHelper.ReloadWebServer()
}

func (repo *MappingCmdRepo) UpdateMarketplaceItem(
	mappingId valueObject.MappingId,
	marketplaceInstalledItemName valueObject.MarketplaceItemName,
) error {
	itemNameStr := marketplaceInstalledItemName.String()
	mappingUpdatedModel := dbModel.Mapping{
		ID:                           mappingId.Uint64(),
		MarketplaceInstalledItemName: &itemNameStr,
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.Mapping{}).
		Where("id = ?", mappingId.Uint64()).
		Updates(&mappingUpdatedModel).Error
}

func (repo *MappingCmdRepo) Update(updateDto dto.UpdateMapping) error {
	err := infraHelper.ValidateWebServerConfig()
	if err != nil {
		return err
	}

	mappingEntity, err := repo.mappingQueryRepo.ReadFirst(dto.ReadMappingsRequest{
		MappingId: &updateDto.Id,
	})
	if err != nil {
		return err
	}

	updateMap := map[string]interface{}{}

	if updateDto.Path != nil {
		updateMap["path"] = updateDto.Path.String()
	}

	if updateDto.MatchPattern != nil {
		updateMap["match_pattern"] = updateDto.MatchPattern.String()
	}

	if updateDto.TargetType != nil {
		updateMap["target_type"] = updateDto.TargetType.String()
	}

	if updateDto.TargetValue != nil {
		targetValueStr := updateDto.TargetValue.String()
		updateMap["target_value"] = &targetValueStr
	}

	if updateDto.TargetHttpResponseCode != nil {
		targetHttpResponseCodeStr := updateDto.TargetHttpResponseCode.String()
		updateMap["target_http_response_code"] = &targetHttpResponseCodeStr
	}

	if updateDto.ShouldUpgradeInsecureRequests != nil {
		updateMap["should_upgrade_insecure_requests"] = updateDto.ShouldUpgradeInsecureRequests
	}

	if updateDto.MappingSecurityRuleId != nil {
		updateMap["mapping_security_rule_id"] = updateDto.MappingSecurityRuleId.Uint64()
	}
	if slices.Contains(updateDto.ClearableFields, "mappingSecurityRuleId") {
		updateMap["mapping_security_rule_id"] = nil
	}

	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.Mapping{}).
		Where("id = ?", updateDto.Id.Uint64()).
		Updates(updateMap).Error
	if err != nil {
		return errors.New("DbUpdateMappingError: " + err.Error())
	}

	err = repo.RecreateMappingFile(mappingEntity.Hostname)
	if err != nil {
		return errors.New("RecreateMappingFileError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}

func (repo *MappingCmdRepo) Delete(mappingId valueObject.MappingId) error {
	err := infraHelper.ValidateWebServerConfig()
	if err != nil {
		return err
	}

	mappingEntity, err := repo.mappingQueryRepo.ReadFirst(dto.ReadMappingsRequest{
		MappingId: &mappingId,
	})
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Delete(dbModel.Mapping{}, mappingId.Uint64()).Error
	if err != nil {
		return err
	}

	err = repo.RecreateMappingFile(mappingEntity.Hostname)
	if err != nil {
		return errors.New("RecreateMappingFileError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}

func (repo *MappingCmdRepo) recreateSecurityRuleFile(
	mappingSecurityRuleId valueObject.MappingSecurityRuleId,
) error {
	ruleEntity, err := repo.mappingQueryRepo.ReadFirstSecurityRule(
		dto.ReadMappingSecurityRulesRequest{MappingSecurityRuleId: &mappingSecurityRuleId},
	)
	if err != nil {
		return errors.New("ReadMappingSecurityRuleError: " + err.Error())
	}

	ruleGlobalTemplateStr := `{{- if and .RpsSoftLimitPerIp .RpsHardLimitPerIp -}}
limit_req_zone $binary_remote_addr zone=rps_limit_{{ .Id }}:10m rate={{ .RpsSoftLimitPerIp }}r/s; #RpsSoftLimitPerIp
{{- else if .RpsSoftLimitPerIp -}}
limit_req_zone $binary_remote_addr zone=rps_limit_{{ .Id }}:10m rate={{ .RpsSoftLimitPerIp }}r/s; #RpsSoftLimitPerIp
{{- else if .RpsHardLimitPerIp -}}
limit_req_zone $binary_remote_addr zone=rps_limit_{{ .Id }}:10m rate={{ .RpsHardLimitPerIp }}r/s; #RpsHardLimitPerIp
{{- end -}}
{{- if .MaxConnectionsPerIp }}
limit_conn_zone $binary_remote_addr zone=conn_limit_{{ .Id }}:10m; #MaxConnectionsPerIp
{{- end }}
`

	ruleGlobalTemplateStrPtr, err := template.New("ruleGlobalFile").Parse(ruleGlobalTemplateStr)
	if err != nil {
		return errors.New("GlobalTemplateParsingError: " + err.Error())
	}

	var ruleGlobalFileContent strings.Builder
	err = ruleGlobalTemplateStrPtr.Execute(&ruleGlobalFileContent, ruleEntity)
	if err != nil {
		return errors.New("GlobalTemplateExecutionError: " + err.Error())
	}

	err = infraHelper.MakeDir(infraEnvs.MappingsSecurityRulesConfDir)
	if err != nil {
		return errors.New("CreateSecurityRulesDirError: " + err.Error())
	}

	ruleGlobalFilePath := infraEnvs.MappingsSecurityRulesConfDir + "/" +
		mappingSecurityRuleId.String() + ".global.conf"
	err = infraHelper.UpdateFile(ruleGlobalFilePath, ruleGlobalFileContent.String(), true)
	if err != nil {
		return errors.New("CreateSecurityRuleGlobalFileError: " + err.Error())
	}

	ruleEmbeddableTemplateStr := `{{- if and .RpsSoftLimitPerIp .RpsHardLimitPerIp -}}
limit_req zone=rps_limit_{{ .Id }} burst={{ .RpsHardLimitPerIp }}; #RpsHardLimitPerIp
{{- else if or .RpsSoftLimitPerIp .RpsHardLimitPerIp -}}
limit_req zone=rps_limit_{{ .Id }}; #RpsLimitPerIp
{{- end -}}
{{- if or .RpsSoftLimitPerIp .RpsHardLimitPerIp -}}
{{- if .ResponseCodeOnMaxRequests }}
limit_req_status {{ .ResponseCodeOnMaxRequests }}; #ResponseCodeOnMaxRequests
{{- end }}
{{- end -}}
{{- if .MaxConnectionsPerIp }}
limit_conn conn_limit_{{ .Id }} {{ .MaxConnectionsPerIp }}; #MaxConnectionsPerIp
{{- if .ResponseCodeOnMaxConnections }}
limit_conn_status {{ .ResponseCodeOnMaxConnections }}; #ResponseCodeOnMaxConnections
{{- end }}
{{- end }}
{{- if .BandwidthBpsLimitPerConnection }}
limit_rate {{ .BandwidthBpsLimitPerConnection }}; #BandwidthBpsLimitPerConnection
{{- if .BandwidthLimitOnlyAfterBytes }}
limit_rate_after {{ .BandwidthLimitOnlyAfterBytes }}; #BandwidthLimitOnlyAfterBytes
{{- end }}
{{- end }}
{{- if .AllowedIps }}
#AllowedIps
{{- range .AllowedIps }}
allow {{ . }};
{{- end }}
{{- end }}
{{- if .BlockedIps }}
#BlockedIps
{{- range .BlockedIps }}
deny {{ . }};
{{- end }}
{{- end }}
`

	ruleEmbeddableTemplateStrPtr, err := template.New("ruleEmbeddableFile").Parse(ruleEmbeddableTemplateStr)
	if err != nil {
		return errors.New("EmbeddableTemplateParsingError: " + err.Error())
	}

	var ruleEmbeddableFileContent strings.Builder
	err = ruleEmbeddableTemplateStrPtr.Execute(&ruleEmbeddableFileContent, ruleEntity)
	if err != nil {
		return errors.New("EmbeddableTemplateExecutionError: " + err.Error())
	}

	ruleEmbeddableFilePath := infraEnvs.MappingsSecurityRulesConfDir + "/" +
		mappingSecurityRuleId.String() + ".embeddable.conf"
	err = infraHelper.UpdateFile(ruleEmbeddableFilePath, ruleEmbeddableFileContent.String(), true)
	if err != nil {
		return errors.New("CreateSecurityRuleEmbeddableFileError: " + err.Error())
	}

	err = infraHelper.UpdateOwnershipForWebServerUse(infraEnvs.MappingsSecurityRulesConfDir, true, false)
	if err != nil {
		return errors.New("UpdateSecurityRulesDirOwnershipError: " + err.Error())
	}

	return nil
}

func (repo *MappingCmdRepo) RecreateSecurityRuleFiles() error {
	responseDto, err := repo.mappingQueryRepo.ReadSecurityRule(
		dto.ReadMappingSecurityRulesRequest{Pagination: dto.PaginationUnpaginated},
	)
	if err != nil {
		return err
	}

	for _, ruleEntity := range responseDto.MappingSecurityRules {
		err = repo.recreateSecurityRuleFile(ruleEntity.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *MappingCmdRepo) CreateSecurityRule(
	createDto dto.CreateMappingSecurityRule,
) (ruleId valueObject.MappingSecurityRuleId, err error) {
	var descriptionPtr *string
	if createDto.Description != nil {
		descriptionStr := createDto.Description.String()
		descriptionPtr = &descriptionStr
	}

	allowedIps := []string{}
	for _, ipAddress := range createDto.AllowedIps {
		allowedIps = append(allowedIps, ipAddress.String())
	}

	blockedIps := []string{}
	for _, ipAddress := range createDto.BlockedIps {
		blockedIps = append(blockedIps, ipAddress.String())
	}

	var bandwidthBpsLimitPerConnectionPtr *uint64
	if createDto.BandwidthBpsLimitPerConnection != nil {
		perConnectionUint64 := createDto.BandwidthBpsLimitPerConnection.Uint64()
		bandwidthBpsLimitPerConnectionPtr = &perConnectionUint64
	}

	var bandwidthLimitOnlyAfterBytesPtr *uint64
	if createDto.BandwidthLimitOnlyAfterBytes != nil {
		afterBytesUint64 := createDto.BandwidthLimitOnlyAfterBytes.Uint64()
		bandwidthLimitOnlyAfterBytesPtr = &afterBytesUint64
	}

	securityRuleModel := dbModel.MappingSecurityRule{
		Name:                           createDto.Name.String(),
		Description:                    descriptionPtr,
		AllowedIps:                     allowedIps,
		BlockedIps:                     blockedIps,
		RpsSoftLimitPerIp:              createDto.RpsSoftLimitPerIp,
		RpsHardLimitPerIp:              createDto.RpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      createDto.ResponseCodeOnMaxRequests,
		MaxConnectionsPerIp:            createDto.MaxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: bandwidthBpsLimitPerConnectionPtr,
		BandwidthLimitOnlyAfterBytes:   bandwidthLimitOnlyAfterBytesPtr,
		ResponseCodeOnMaxConnections:   createDto.ResponseCodeOnMaxConnections,
	}

	err = repo.persistentDbSvc.Handler.Create(&securityRuleModel).Error
	if err != nil {
		return ruleId, errors.New("CreateMappingSecurityRuleInfraError")
	}

	ruleId, err = valueObject.NewMappingSecurityRuleId(securityRuleModel.ID)
	if err != nil {
		return ruleId, err
	}

	err = repo.recreateSecurityRuleFile(ruleId)
	if err != nil {
		return ruleId, errors.New("RecreateMappingSecurityRuleFileError: " + err.Error())
	}

	err = infraHelper.ValidateWebServerConfig()
	if err != nil {
		err = repo.persistentDbSvc.Handler.Delete(&securityRuleModel).Error
		if err != nil {
			return ruleId, err
		}

		err = repo.recreateSecurityRuleFile(ruleId)
		if err != nil {
			return ruleId, errors.New("RecreateMappingSecurityRuleFileError: " + err.Error())
		}
	}

	return ruleId, infraHelper.ReloadWebServer()
}

func (repo *MappingCmdRepo) UpdateSecurityRule(
	updateDto dto.UpdateMappingSecurityRule,
) error {
	updateMap := map[string]interface{}{}

	if updateDto.Name != nil {
		updateMap["name"] = updateDto.Name.String()
	}

	if updateDto.Description != nil {
		updateMap["description"] = updateDto.Description.String()
	}
	if slices.Contains(updateDto.ClearableFields, "description") {
		updateMap["description"] = nil
	}

	if updateDto.RpsSoftLimitPerIp != nil {
		updateMap["rps_soft_limit_per_ip"] = *updateDto.RpsSoftLimitPerIp
	}

	if updateDto.RpsHardLimitPerIp != nil {
		updateMap["rps_hard_limit_per_ip"] = *updateDto.RpsHardLimitPerIp
	}

	if updateDto.ResponseCodeOnMaxRequests != nil {
		updateMap["response_code_on_max_requests"] = *updateDto.ResponseCodeOnMaxRequests
	}

	if updateDto.MaxConnectionsPerIp != nil {
		updateMap["max_connections_per_ip"] = *updateDto.MaxConnectionsPerIp
	}

	if updateDto.BandwidthBpsLimitPerConnection != nil {
		updateMap["bandwidth_bps_limit_per_connection"] = updateDto.BandwidthBpsLimitPerConnection.Uint64()
	}

	if updateDto.BandwidthLimitOnlyAfterBytes != nil {
		updateMap["bandwidth_limit_only_after_bytes"] = updateDto.BandwidthLimitOnlyAfterBytes.Uint64()
	}

	if updateDto.ResponseCodeOnMaxConnections != nil {
		updateMap["response_code_on_max_connections"] = *updateDto.ResponseCodeOnMaxConnections
	}

	err := repo.persistentDbSvc.Handler.Model(&dbModel.MappingSecurityRule{}).
		Where("id = ?", updateDto.Id.Uint64()).Updates(updateMap).Error
	if err != nil {
		return errors.New("UpdateMappingSecurityRuleInfraError: " + err.Error())
	}

	updateModel := dbModel.MappingSecurityRule{}
	allowedIps := []string{}
	for _, ipAddress := range updateDto.AllowedIps {
		allowedIps = append(allowedIps, ipAddress.String())
	}
	if len(allowedIps) > 0 {
		updateModel.AllowedIps = allowedIps
	}

	blockedIps := []string{}
	for _, ipAddress := range updateDto.BlockedIps {
		blockedIps = append(blockedIps, ipAddress.String())
	}
	if len(blockedIps) > 0 {
		updateModel.BlockedIps = blockedIps
	}

	if len(updateModel.AllowedIps) > 0 || len(updateModel.BlockedIps) > 0 {
		err = repo.persistentDbSvc.Handler.Model(&dbModel.MappingSecurityRule{}).
			Where("id = ?", updateDto.Id.Uint64()).Updates(updateModel).Error
		if err != nil {
			return errors.New("UpdateMappingSecurityRuleInfraError: " + err.Error())
		}
	}

	err = repo.recreateSecurityRuleFile(updateDto.Id)
	if err != nil {
		return errors.New("RecreateMappingSecurityRuleFileError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}

func (repo *MappingCmdRepo) DeleteSecurityRule(
	ruleId valueObject.MappingSecurityRuleId,
) error {
	err := infraHelper.ValidateWebServerConfig()
	if err != nil {
		return err
	}

	ruleEntity, err := repo.mappingQueryRepo.ReadFirstSecurityRule(
		dto.ReadMappingSecurityRulesRequest{MappingSecurityRuleId: &ruleId},
	)
	if err != nil {
		return errors.New("ReadMappingSecurityRuleError: " + err.Error())
	}

	err = repo.persistentDbSvc.Handler.Delete(
		dbModel.MappingSecurityRule{}, ruleEntity.Id.Uint64(),
	).Error
	if err != nil {
		return err
	}

	toRemoveFilePaths := []string{
		infraEnvs.MappingsSecurityRulesConfDir + "/" + ruleEntity.Id.String() + ".global.conf",
		infraEnvs.MappingsSecurityRulesConfDir + "/" + ruleEntity.Id.String() + ".embeddable.conf",
	}
	for _, toRemoveFilePath := range toRemoveFilePaths {
		err = os.Remove(toRemoveFilePath)
		if err != nil {
			return errors.New("RemoveSecurityRuleFileError: " + err.Error())
		}
	}

	return infraHelper.ReloadWebServer()
}
