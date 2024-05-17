package vhostInfra

import (
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	infraData "github.com/speedianet/os/src/infra/infraData"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type VirtualHostCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	queryRepo       *VirtualHostQueryRepo
}

func NewVirtualHostCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostCmdRepo {
	vhostQueryRepo := NewVirtualHostQueryRepo(persistentDbSvc)

	return &VirtualHostCmdRepo{
		persistentDbSvc: persistentDbSvc,
		queryRepo:       vhostQueryRepo,
	}
}

func (repo *VirtualHostCmdRepo) webServerFileFactory(
	hostname valueObject.Fqdn,
	aliases []valueObject.Fqdn,
	publicDir valueObject.UnixFilePath,
	mappingFilePath valueObject.UnixFilePath,
) (string, error) {
	hostnameStr := hostname.String()

	valuesToInterpolate := map[string]interface{}{
		"Hostname":        hostname,
		"Aliases":         aliases,
		"PublicDirectory": publicDir,
		"CertPath":        infraData.GlobalConfigs.PkiConfDir + "/" + hostnameStr + ".crt",
		"KeyPath":         infraData.GlobalConfigs.PkiConfDir + "/" + hostnameStr + ".key",
		"MappingFilePath": mappingFilePath,
	}

	webServerConfigTemplate := `server {
    listen 80;
    listen 443 ssl;
    server_name {{ .Hostname }} www.{{ .Hostname }}{{ range .Aliases }} {{ .String }} www.{{ .String }}{{ end }};

    root {{ .PublicDirectory }};

    ssl_certificate {{ .CertPath }};
    ssl_certificate_key {{ .KeyPath }};

    access_log /app/logs/nginx/{{ .Hostname }}_access.log combined buffer=512k flush=1m;
    error_log /app/logs/nginx/{{ .Hostname }}_error.log warn;

    include /etc/nginx/std.conf;
    include {{ .MappingFilePath }};
}`

	webServerConfigTemplatePtr, err := template.
		New("webServerConfigFile").
		Parse(webServerConfigTemplate)
	if err != nil {
		return "", errors.New("TemplateParsingError: " + err.Error())
	}

	var webServerConfigFileContent strings.Builder
	err = webServerConfigTemplatePtr.Execute(
		&webServerConfigFileContent,
		valuesToInterpolate,
	)
	if err != nil {
		return "", errors.New("TemplateExecutionError: " + err.Error())
	}

	return webServerConfigFileContent.String(), nil
}

func (repo *VirtualHostCmdRepo) createWebServerFile(
	publicDir valueObject.UnixFilePath,
	hostname valueObject.Fqdn,
) error {
	aliases, err := repo.queryRepo.ReadAliasesByHostname(hostname)
	if err != nil {
		return errors.New("GetAliasesByHostnameError: " + err.Error())
	}

	aliasesHostnames := []valueObject.Fqdn{}
	for _, alias := range aliases {
		aliasesHostnames = append(aliasesHostnames, alias.Hostname)
	}

	hostnameStr := hostname.String()
	if infraHelper.IsPrimaryVirtualHost(hostname) {
		hostnameStr = "primary"
	}

	mappingFilePathStr := infraData.GlobalConfigs.MappingsConfDir + "/" + hostnameStr + ".conf"
	mappingFilePath, err := valueObject.NewUnixFilePath(mappingFilePathStr)
	if err != nil {
		return errors.New(err.Error() + ": " + mappingFilePathStr)
	}
	err = infraHelper.UpdateFile(mappingFilePath.String(), "", false)
	if err != nil {
		return errors.New("CreateMappingFileFailed")
	}

	webServerConfigFileContent, err := repo.webServerFileFactory(
		hostname,
		aliasesHostnames,
		publicDir,
		mappingFilePath,
	)
	if err != nil {
		return err
	}

	webServerFilePathStr := infraData.GlobalConfigs.VirtualHostsConfDir + "/" + hostnameStr + ".conf"
	webServerFilePath, err := valueObject.NewUnixFilePath(webServerFilePathStr)
	if err != nil {
		return errors.New(err.Error() + ": " + webServerFilePathStr)
	}
	err = infraHelper.UpdateFile(
		webServerFilePath.String(),
		webServerConfigFileContent,
		true,
	)
	if err != nil {
		return errors.New("CreateNginxConfFileFailed")
	}

	return infraHelper.ReloadWebServer()
}

func (repo *VirtualHostCmdRepo) persistVirtualHost(
	createDto dto.CreateVirtualHost,
	publicDir valueObject.UnixFilePath,
) error {
	var parentHostnamePtr *string
	if createDto.ParentHostname != nil {
		parentHostnameStr := createDto.ParentHostname.String()
		parentHostnamePtr = &parentHostnameStr
	}

	model := dbModel.VirtualHost{
		Hostname:       createDto.Hostname.String(),
		Type:           createDto.Type.String(),
		RootDirectory:  publicDir.String(),
		ParentHostname: parentHostnamePtr,
	}
	return repo.persistentDbSvc.Handler.Create(&model).Error
}

func (repo *VirtualHostCmdRepo) createAlias(createDto dto.CreateVirtualHost) error {
	parentVhost, err := repo.queryRepo.ReadByHostname(*createDto.ParentHostname)
	if err != nil {
		return errors.New("GetParentVhostError: " + err.Error())
	}

	err = repo.persistVirtualHost(createDto, parentVhost.RootDirectory)
	if err != nil {
		return err
	}

	return repo.createWebServerFile(
		parentVhost.RootDirectory,
		parentVhost.Hostname,
	)
}

func (repo *VirtualHostCmdRepo) updateDirsOwnership(
	publicDir valueObject.UnixFilePath,
) error {
	directories := []string{
		publicDir.String(),
		"/app/conf/nginx",
		infraData.GlobalConfigs.PkiConfDir,
	}

	for _, directory := range directories {
		_, err := infraHelper.RunCmd(
			"chown",
			"-R",
			"nobody:nogroup",
			directory,
		)
		if err != nil {
			return errors.New("ChownNecessaryDirectoriesFailed")
		}
	}

	return nil
}

func (repo *VirtualHostCmdRepo) Create(createDto dto.CreateVirtualHost) error {
	if createDto.Type.String() == "alias" {
		return repo.createAlias(createDto)
	}

	hostnameStr := createDto.Hostname.String()

	publicDirStr := infraData.GlobalConfigs.PrimaryPublicDir + "/" + hostnameStr
	publicDir, err := valueObject.NewUnixFilePath(publicDirStr)
	if err != nil {
		return errors.New(err.Error() + ": " + publicDirStr)
	}

	err = infraHelper.MakeDir(publicDirStr)
	if err != nil {
		return errors.New("MakePublicHtmlDirFailed")
	}

	err = infraHelper.CreateSelfSignedSsl(
		infraData.GlobalConfigs.PkiConfDir,
		hostnameStr,
	)
	if err != nil {
		return errors.New("GenerateSelfSignedCertFailed")
	}

	err = repo.updateDirsOwnership(publicDir)
	if err != nil {
		return err
	}

	err = repo.persistVirtualHost(createDto, publicDir)
	if err != nil {
		return err
	}

	return repo.createWebServerFile(publicDir, createDto.Hostname)
}

func (repo *VirtualHostCmdRepo) deleteWebServerFile(
	hostname valueObject.Fqdn,
) error {
	hostnameStr := hostname.String()
	if infraHelper.IsPrimaryVirtualHost(hostname) {
		hostnameStr = "primary"
	}

	mappingFilePathStr := infraData.GlobalConfigs.MappingsConfDir + "/" + hostnameStr + ".conf"
	err := os.Remove(mappingFilePathStr)
	if err != nil {
		return err
	}

	webServerFilePathStr := infraData.GlobalConfigs.VirtualHostsConfDir + "/" + hostnameStr + ".conf"
	err = os.Remove(webServerFilePathStr)
	if err != nil {
		return err
	}

	return infraHelper.ReloadWebServer()
}

func (repo *VirtualHostCmdRepo) Delete(vhost entity.VirtualHost) error {
	err := repo.persistentDbSvc.Handler.
		Where(
			"Id = ? OR ParentHostname = ?",
			vhost.Id.Get(),
			vhost.Hostname.String(),
		).
		Delete(dbModel.VirtualHost{}).Error
	if err != nil {
		return err
	}

	if vhost.Type.String() == "alias" {
		parentVhost, err := repo.queryRepo.ReadByHostname(*vhost.ParentHostname)
		if err != nil {
			return errors.New("GetParentVhost: " + err.Error())
		}
		vhost = parentVhost

		return repo.createWebServerFile(
			vhost.RootDirectory,
			vhost.Hostname,
		)
	}

	err = repo.deleteWebServerFile(vhost.Hostname)
	if err != nil {
		return errors.New("DeleteWebServerFileError: " + err.Error())
	}

	return nil
}
