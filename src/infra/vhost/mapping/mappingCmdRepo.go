package mappingInfra

import (
	"errors"
	"log"
	"strings"
	"text/template"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

type MappingCmdRepo struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	mappingQueryRepo *MappingQueryRepo
	vhostCmdRepo     vhostInfra.VirtualHostCmdRepo
}

func NewMappingCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MappingCmdRepo {
	mappingQueryRepo := NewMappingQueryRepo(persistentDbSvc)
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	return &MappingCmdRepo{
		persistentDbSvc:  persistentDbSvc,
		mappingQueryRepo: mappingQueryRepo,
		vhostCmdRepo:     vhostCmdRepo,
	}
}

func (repo *MappingCmdRepo) getServiceMappingConfig(svcNameStr string) (string, error) {
	svcMappingConfig := ""

	serviceName, err := valueObject.NewServiceName(svcNameStr)
	if err != nil {
		return "", errors.New(err.Error() + ": " + svcNameStr)
	}

	svcQueryRepo := servicesInfra.ServicesQueryRepo{}
	service, err := svcQueryRepo.GetByName(serviceName)
	if err != nil {
		return "", errors.New("GetServiceByNameError")
	}

	protocolPortsMap := map[string]string{}
	for _, svcPortBinding := range service.PortBindings {
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
) string {
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

	locationUri := pathStr
	if modifier != "" {
		locationUri = modifier + " " + pathStr
	}

	return locationUri
}

func (repo *MappingCmdRepo) recreateMappingFile(
	mappingHostname valueObject.Fqdn,
) error {
	mappings, err := repo.mappingQueryRepo.GetByHostname(mappingHostname)
	if err != nil {
		return err
	}

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	mappingFilePath, err := vhostQueryRepo.GetVirtualHostMappingsFilePath(
		mappingHostname,
	)
	if err != nil {
		return errors.New("GetVirtualHostMappingsFilePathError: " + err.Error())
	}

	mappingConfigTemplate := `{{- range . -}}
location {{ parseLocationUri .MatchPattern .Path }} {
	{{- if eq .TargetType "url" "response-code" }}
	return {{ .TargetHttpResponseCode }} {{ .TargetValue }};
	{{- end }}
	{{- if eq .TargetType "service" }}
{{ getServiceMappingConfig .TargetValue.String }}
	{{- end }}
	{{- if eq .TargetType "inline-html" }}
	add_header Content-Type text/html;
	return {{ .TargetHttpResponseCode }} "{{ .TargetValue }}";
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
	err = mappingTemplatePtr.Execute(&mappingFileContent, mappings)
	if err != nil {
		return errors.New("TemplateExecutionError: " + err.Error())
	}

	shouldOverwrite := true
	return infraHelper.UpdateFile(
		mappingFilePath.String(),
		mappingFileContent.String(),
		shouldOverwrite,
	)
}

func (repo *MappingCmdRepo) Create(
	createDto dto.CreateMapping,
) (valueObject.MappingId, error) {
	var mappingId valueObject.MappingId

	isServiceMapping := createDto.TargetType.String() == "service"
	isPhpServiceMapping := isServiceMapping && createDto.TargetValue.String() == "php"
	if isPhpServiceMapping {
		err := repo.vhostCmdRepo.CreatePhpVirtualHost(createDto.Hostname)
		if err != nil {
			return mappingId, err
		}
	}

	mappingModel := dbModel.Mapping{}.AddDtoToModel(createDto)
	createResult := repo.persistentDbSvc.Handler.Create(&mappingModel)
	if createResult.Error != nil {
		return mappingId, createResult.Error
	}
	mappingId, err := valueObject.NewMappingId(mappingModel.ID)
	if err != nil {
		return mappingId, err
	}

	err = repo.recreateMappingFile(createDto.Hostname)
	if err != nil {
		log.Printf("NewRecreateMappingFileFunctionError: %s", err.Error())
	}

	return mappingId, repo.vhostCmdRepo.ReloadWebServer()
}

func (repo *MappingCmdRepo) Delete(mappingId valueObject.MappingId) error {
	mapping, err := repo.mappingQueryRepo.GetById(mappingId)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.Delete(
		dbModel.Mapping{},
		mappingId.Get(),
	).Error
	if err != nil {
		return err
	}

	err = repo.recreateMappingFile(mapping.Hostname)
	if err != nil {
		return err
	}

	return repo.vhostCmdRepo.ReloadWebServer()
}

func (repo *MappingCmdRepo) DeleteAuto(
	serviceName valueObject.ServiceName,
) error {
	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return errors.New("PrimaryVhostNotFound")
	}

	primaryVhostMappings, err := repo.mappingQueryRepo.GetByHostname(primaryVhost)
	if err != nil {
		return errors.New("GetPrimaryVhostMappingsError: " + err.Error())
	}

	var mappingIdToDelete *valueObject.MappingId
	for _, primaryVhostMapping := range primaryVhostMappings {
		if primaryVhostMapping.TargetType.String() != "service" {
			continue
		}

		if primaryVhostMapping.TargetValue.String() != serviceName.String() {
			continue
		}

		mappingIdToDelete = &primaryVhostMapping.Id
	}

	if mappingIdToDelete == nil {
		return nil
	}

	return repo.Delete(*mappingIdToDelete)
}

func (repo *MappingCmdRepo) Recreate(mapping entity.Mapping) error {
	err := repo.Delete(mapping.Id)
	if err != nil {
		return err
	}

	createDto := dto.NewCreateMapping(
		mapping.Hostname,
		mapping.Path,
		mapping.MatchPattern,
		mapping.TargetType,
		mapping.TargetValue,
		mapping.TargetHttpResponseCode,
	)

	_, err = repo.Create(createDto)
	return err
}
