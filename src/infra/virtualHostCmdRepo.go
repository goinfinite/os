package infra

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

type VirtualHostCmdRepo struct {
}

func (repo VirtualHostCmdRepo) reloadWebServer() error {
	_, err := infraHelper.RunCmd(
		"nginx",
		"-t",
	)
	if err != nil {
		return errors.New("NginxConfigTestFailed")
	}

	_, err = infraHelper.RunCmd(
		"nginx",
		"-s",
		"reload",
	)
	if err != nil {
		return errors.New("NginxReloadFailed")
	}

	return nil
}

func (repo VirtualHostCmdRepo) getAliasConfigFile(
	parentHostname valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	vhostFileStr := "/app/conf/nginx/" + parentHostname.String() + ".conf"

	isParentPrimaryDomain := VirtualHostQueryRepo{}.IsVirtualHostPrimaryDomain(
		parentHostname,
	)
	if isParentPrimaryDomain {
		vhostFileStr = "/app/conf/nginx/primary.conf"
	}

	return valueObject.UnixFilePath(vhostFileStr), nil
}

func (repo VirtualHostCmdRepo) addAlias(addDto dto.AddVirtualHost) error {
	vhostFile, err := repo.getAliasConfigFile(*addDto.ParentHostname)
	if err != nil {
		return errors.New("GetAliasConfigFileFailed")
	}
	vhostFileStr := vhostFile.String()

	hostnameStr := addDto.Hostname.String()

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		`/server_name/ s/;$/ `+hostnameStr+` www.`+hostnameStr+`;/`,
		vhostFileStr,
	)
	if err != nil {
		return errors.New("AddAliasFailed")
	}

	// TODO: Regenerate cert for primary domain to include new alias

	return repo.reloadWebServer()
}

func (repo VirtualHostCmdRepo) Add(addDto dto.AddVirtualHost) error {
	hostnameStr := addDto.Hostname.String()

	if addDto.Type.String() == "alias" {
		return repo.addAlias(addDto)
	}

	publicDir := "/app/html/" + hostnameStr
	certPath := "/app/conf/pki/" + hostnameStr + ".crt"
	keyPath := "/app/conf/pki/" + hostnameStr + ".key"
	mappingFilePath := "/app/conf/nginx/mapping/" + hostnameStr + ".conf"

	nginxConf := `server {
    listen 80;
    listen 443 ssl;
    server_name ` + hostnameStr + ` www.` + hostnameStr + `;

    root ` + publicDir + `;

    ssl_certificate ` + certPath + `;
    ssl_certificate_key ` + keyPath + `;

    access_log /app/logs/nginx/` + hostnameStr + `_access.log combined buffer=512k flush=1m;
    error_log /app/logs/nginx/` + hostnameStr + `_error.log warn;

    include /etc/nginx/std.conf;
    include ` + mappingFilePath + `;
}
`
	err := infraHelper.UpdateFile(
		"/app/conf/nginx/"+hostnameStr+".conf",
		nginxConf,
		true,
	)
	if err != nil {
		return errors.New("CreateNginxConfFileFailed")
	}

	err = infraHelper.UpdateFile(
		mappingFilePath,
		"",
		true,
	)
	if err != nil {
		return errors.New("CreateMappingFileFailed")
	}

	err = infraHelper.MakeDir(publicDir)
	if err != nil {
		return errors.New("MakePublicHtmlDirFailed")
	}

	_, err = infraHelper.RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		keyPath,
		"-out",
		certPath,
		"-subj",
		"/C=US/ST=California/L=LosAngeles/O=Acme/CN="+hostnameStr,
	)
	if err != nil {
		return errors.New("GenerateSelfSignedCertFailed")
	}

	directories := []string{
		publicDir,
		"/app/conf/nginx",
		"/app/conf/pki",
	}
	for _, directory := range directories {
		_, err = infraHelper.RunCmd(
			"chown",
			"-R",
			"nobody:nogroup",
			directory,
		)
		if err != nil {
			return errors.New("ChownNecessaryDirectoriesFailed")
		}
	}

	return repo.reloadWebServer()
}

func (repo VirtualHostCmdRepo) deleteAlias(vhost entity.VirtualHost) error {
	vhostFile, err := repo.getAliasConfigFile(*vhost.ParentHostname)
	if err != nil {
		return errors.New("GetAliasConfigFileFailed")
	}
	vhostFileStr := vhostFile.String()

	hostnameStr := vhost.Hostname.String()

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		`/server_name/ s/ `+hostnameStr+` www.`+hostnameStr+`//`,
		vhostFileStr,
	)
	if err != nil {
		return errors.New("DeleteAliasFailed")
	}

	return repo.reloadWebServer()
}

func (repo VirtualHostCmdRepo) Delete(vhost entity.VirtualHost) error {
	hostnameStr := vhost.Hostname.String()
	if vhost.Type.String() == "alias" {
		return repo.deleteAlias(vhost)
	}

	_, err := infraHelper.RunCmd(
		"rm",
		"-rf",
		"/app/html/"+hostnameStr,
		"/app/conf/nginx/"+hostnameStr+".conf",
		"/app/conf/pki/"+hostnameStr+".crt",
		"/app/conf/pki/"+hostnameStr+".key",
		"/app/conf/nginx/mapping/"+hostnameStr+".conf",
	)
	if err != nil {
		return errors.New("DeleteVirtualHostFailed")
	}

	return repo.reloadWebServer()
}

func (repo VirtualHostCmdRepo) mappingToLocationUri(
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

func (repo VirtualHostCmdRepo) serviceLocationContentFactory(
	addMapping dto.AddMapping,
) (string, error) {
	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	serviceEntity, err := servicesQueryRepo.GetByName(*addMapping.TargetServiceName)
	if err != nil {
		return "", errors.New("GetServiceByNameFailed")
	}

	protocolPortsMap := map[string]string{}
	for _, svcPortBinding := range serviceEntity.PortBindings {
		protocolPortsMap[svcPortBinding.Protocol.String()] = svcPortBinding.Port.String()
	}

	locationContent := ""
	isHttpSupported := protocolPortsMap["http"] != ""
	if isHttpSupported {
		locationContent += `
	set $protocol "http";
	set $backend "localhost:` + protocolPortsMap["http"] + `";
`
	}

	isHttpsSupported := protocolPortsMap["https"] != ""
	if isHttpsSupported {
		locationContent += `
	set $protocol "https";
	set $backend "localhost:` + protocolPortsMap["https"] + `";
`
	}

	if isHttpSupported && isHttpsSupported {
		locationContent = `
	set $protocol "http";
	set $backend "localhost:` + protocolPortsMap["http"] + `";

	if ($scheme = https) {
		set $protocol "https";
		set $backend "localhost:` + protocolPortsMap["https"] + `";
	}
`
	}

	isHttpOrHttpsSupported := isHttpSupported || isHttpsSupported

	isGrpcSupported := protocolPortsMap["grpc"] != ""
	if isGrpcSupported && !isHttpOrHttpsSupported {
		locationContent += `
	set $protocol "grpc";
	set $backend "localhost:` + protocolPortsMap["grpc"] + `";
`
	}

	if isGrpcSupported && isHttpOrHttpsSupported {
		locationContent += `
	if ($scheme = grpc) {
		set $protocol "grpc";
		set $backend "localhost:` + protocolPortsMap["grpc"] + `";
	}
`
	}

	isGrpcsSupported := protocolPortsMap["grpcs"] != ""
	if isGrpcsSupported && !isHttpOrHttpsSupported {
		locationContent += `
	set $protocol "grpcs";
	set $backend "localhost:` + protocolPortsMap["grpcs"] + `";
`
	}

	if isGrpcsSupported && isHttpOrHttpsSupported {
		locationContent += `
	if ($scheme = grpcs) {
		set $protocol "grpcs";
		set $backend "localhost:` + protocolPortsMap["grpcs"] + `";
	}
		`
	}

	if isGrpcSupported && !isGrpcsSupported && isHttpOrHttpsSupported {
		locationContent += `
	grpc_set_header Host $host;
	if ($protocol = grpc) {	
		grpc_pass $protocol://$backend;
	}
`
	}

	if !isGrpcSupported && isGrpcsSupported && isHttpOrHttpsSupported {
		locationContent += `
	grpc_set_header Host $host;
	if ($protocol = grpcs) {	
		grpc_pass $protocol://$backend;
	}
`
	}

	isGrpcAndGrpcsSupported := isGrpcSupported && isGrpcsSupported
	if isGrpcAndGrpcsSupported && !isHttpOrHttpsSupported {
		locationContent = `
	set $protocol "grpc";
	set $backend "localhost:` + protocolPortsMap["grpc"] + `";

	if ($scheme = grpcs) {
		set $protocol "grpcs";
		set $backend "localhost:` + protocolPortsMap["grpcs"] + `";
	}
`
	}

	if isGrpcAndGrpcsSupported && isHttpOrHttpsSupported {
		locationContent += `
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
		locationContent += `
	grpc_set_header Host $host;
	grpc_pass $protocol://$backend;
`
	}

	isWsSupported := protocolPortsMap["ws"] != ""
	if isWsSupported && !isHttpOrHttpsSupported {
		locationContent += `
	set $protocol "http";
	set $backend "localhost:` + protocolPortsMap["ws"] + `";
`
	}

	isWssSupported := protocolPortsMap["wss"] != ""
	if !isWsSupported && isWssSupported && !isHttpOrHttpsSupported {
		locationContent += `
	set $protocol "https";
	set $backend "localhost:` + protocolPortsMap["wss"] + `";
`
	}

	isWsAndWssSupported := isWsSupported && isWssSupported
	if isWsAndWssSupported && !isHttpOrHttpsSupported {
		locationContent = `
	set $protocol "http";
	set $backend "localhost:` + protocolPortsMap["ws"] + `";

	if ($scheme = https) {
		set $protocol "https";
		set $backend "localhost:` + protocolPortsMap["wss"] + `";
	}

	proxy_pass $protocol://$backend;
	proxy_set_header Host $host;
`
	}

	isWsOrWssSupported := isWsSupported || isWssSupported
	if isWsOrWssSupported {
		locationContent += `
	proxy_http_version 1.1;
	proxy_set_header Upgrade $http_upgrade;
	proxy_set_header Connection "Upgrade";
`
	}

	if isHttpOrHttpsSupported {
		locationContent += `
	proxy_pass $protocol://$backend;
	proxy_set_header Host $host;
`
	}

	return locationContent, nil
}

func (repo VirtualHostCmdRepo) AddMapping(addMapping dto.AddMapping) error {
	locationUri := repo.mappingToLocationUri(addMapping.MatchPattern, addMapping.Path)

	responseCodeStr := ""
	if addMapping.TargetHttpResponseCode != nil {
		responseCodeStr = addMapping.TargetHttpResponseCode.String()
	}

	locationContent := "	return " + responseCodeStr
	if addMapping.TargetType.String() == "url" {
		locationContent += " " + addMapping.TargetUrl.String()
	}
	locationContent += ";"

	if addMapping.TargetType.String() == "service" {
		var err error
		locationContent, err = repo.serviceLocationContentFactory(addMapping)
		if err != nil {
			return errors.New("ServiceLocationContentFactoryFailed")
		}
	}

	locationBlock := `location ` + locationUri + ` {
` + locationContent + `
}
`

	vhostQueryRepo := VirtualHostQueryRepo{}
	mappingFilePath, err := vhostQueryRepo.GetVirtualHostMappingsFilePath(
		addMapping.Hostname,
	)
	if err != nil {
		return errors.New("GetVirtualHostMappingsFilePathFailed")
	}

	err = infraHelper.UpdateFile(
		mappingFilePath.String(),
		locationBlock,
		false,
	)
	if err != nil {
		return errors.New("AddMappingFailed")
	}

	return repo.reloadWebServer()
}

func (repo VirtualHostCmdRepo) DeleteMapping(mapping entity.Mapping) error {
	locationUri := repo.mappingToLocationUri(mapping.MatchPattern, mapping.Path)

	vhostQueryRepo := VirtualHostQueryRepo{}
	mappingFilePath, err := vhostQueryRepo.GetVirtualHostMappingsFilePath(
		mapping.Hostname,
	)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		`\#location `+locationUri+` {#,/}/d`,
		mappingFilePath.String(),
	)
	if err != nil {
		return err
	}

	return repo.reloadWebServer()
}
