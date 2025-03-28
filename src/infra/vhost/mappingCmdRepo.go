package vhostInfra

import (
	"errors"
	"strings"
	"text/template"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
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

func (repo *MappingCmdRepo) getServiceMappingConfig(
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

func (repo *MappingCmdRepo) parseLocationUri(
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

func (repo *MappingCmdRepo) recreateMappingFile(
	mappingHostname valueObject.Fqdn,
) error {
	mappingsReadResponse, err := repo.mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Pagination: dto.PaginationUnpaginated,
		Hostname:   &mappingHostname,
	})
	if err != nil {
		return err
	}

	vhostQueryRepo := NewVirtualHostQueryRepo(repo.persistentDbSvc)
	mappingFilePath, err := vhostQueryRepo.ReadVirtualHostMappingsFilePath(
		mappingHostname,
	)
	if err != nil {
		return errors.New("GetVirtualHostMappingsFilePathError: " + err.Error())
	}

	mappingConfigTemplate := `{{- range . -}}
location {{ parseLocationUri .MatchPattern .Path }} {
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
{{ getServiceMappingConfig .TargetValue.String }}
	{{- end }}
	{{- if eq .TargetType "static-files" }}
	try_files $uri $uri/ index.html?$query_string;
	{{- end }}
}
{{ end }}`

	mappingTemplatePtr := template.New("mappingFile")
	mappingTemplatePtr = mappingTemplatePtr.Funcs(
		template.FuncMap{
			"parseLocationUri":        repo.parseLocationUri,
			"getServiceMappingConfig": repo.getServiceMappingConfig,
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

	mappingModel := dbModel.NewMapping(
		0, createDto.Hostname.String(), createDto.Path.String(),
		createDto.MatchPattern.String(), createDto.TargetType.String(),
		targetValuePtr, targetHttpResponseCodePtr,
	)

	err = repo.persistentDbSvc.Handler.Create(&mappingModel).Error
	if err != nil {
		return mappingId, err
	}

	mappingId, err = valueObject.NewMappingId(mappingModel.ID)
	if err != nil {
		return mappingId, err
	}

	err = repo.recreateMappingFile(createDto.Hostname)
	if err != nil {
		return mappingId, errors.New("RecreateMappingFileError: " + err.Error())
	}

	err = infraHelper.ValidateWebServerConfig()
	if err != nil {
		err = repo.persistentDbSvc.Handler.Delete(&mappingModel).Error
		if err != nil {
			return mappingId, err
		}

		err = repo.recreateMappingFile(createDto.Hostname)
		if err != nil {
			return mappingId, errors.New("RecreateMappingFileError: " + err.Error())
		}
	}

	return mappingId, infraHelper.ReloadWebServer()
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

	err = repo.recreateMappingFile(mappingEntity.Hostname)
	if err != nil {
		return errors.New("RecreateMappingFileError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}
