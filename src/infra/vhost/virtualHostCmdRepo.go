package vhostInfra

import (
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
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

func (repo *VirtualHostCmdRepo) webServerUnitFileFactory(
	vhostName valueObject.Fqdn,
	aliasesHostnames []valueObject.Fqdn,
	publicDir valueObject.UnixFilePath,
	mappingFilePath valueObject.UnixFilePath,
) (string, error) {
	vhostNameStr := vhostName.String()

	aliasesHostnamesStr := []string{}
	for _, aliasHostname := range aliasesHostnames {
		aliasesHostnamesStr = append(aliasesHostnamesStr, aliasHostname.String())
	}

	valuesToInterpolate := map[string]interface{}{
		"VhostName":        vhostNameStr,
		"AliasesHostnames": aliasesHostnamesStr,
		"PublicDirectory":  publicDir,
		"CertPath":         infraEnvs.PkiConfDir + "/" + vhostNameStr + ".crt",
		"KeyPath":          infraEnvs.PkiConfDir + "/" + vhostNameStr + ".key",
		"MappingFilePath":  mappingFilePath,
	}

	unitConfTemplate := `server {
    listen 80;
    listen 443 ssl;
    server_name {{ .VhostName }} www.{{ .VhostName }}{{ range $aliasHostname := .AliasesHostnames }} {{ $aliasHostname }} www.{{ $aliasHostname }}{{ end }};

    root {{ .PublicDirectory }};

    ssl_certificate {{ .CertPath }};
    ssl_certificate_key {{ .KeyPath }};

    access_log /app/logs/nginx/{{ .VhostName }}_access.log combined buffer=512k flush=1m;
    error_log /app/logs/nginx/{{ .VhostName }}_error.log warn;

    include /etc/nginx/std.conf;
    include {{ .MappingFilePath }};
}`

	unitConfTemplatePtr, err := template.
		New("webServerConfUnitFile").
		Parse(unitConfTemplate)
	if err != nil {
		return "", errors.New("TemplateParsingError: " + err.Error())
	}

	var unitConfFileContent strings.Builder
	err = unitConfTemplatePtr.Execute(
		&unitConfFileContent,
		valuesToInterpolate,
	)
	if err != nil {
		return "", errors.New("TemplateExecutionError: " + err.Error())
	}

	return unitConfFileContent.String(), nil
}

func (repo *VirtualHostCmdRepo) createWebServerUnitFile(
	vhostName valueObject.Fqdn,
	publicDir valueObject.UnixFilePath,
) error {
	aliases, err := repo.queryRepo.ReadAliasesByParentHostname(vhostName)
	if err != nil {
		return errors.New("GetAliasesByHostnameError: " + err.Error())
	}

	aliasesHostnames := []valueObject.Fqdn{}
	for _, alias := range aliases {
		aliasesHostnames = append(aliasesHostnames, alias.Hostname)
	}

	vhostFileNameStr := vhostName.String() + ".conf"
	if infraHelper.IsPrimaryVirtualHost(vhostName) {
		vhostFileNameStr = infraEnvs.PrimaryVhostFileName
	}

	mappingFilePathStr := infraEnvs.MappingsConfDir + "/" + vhostFileNameStr
	mappingFilePath, err := valueObject.NewUnixFilePath(mappingFilePathStr)
	if err != nil {
		return errors.New(err.Error() + ": " + mappingFilePathStr)
	}
	err = infraHelper.UpdateFile(mappingFilePath.String(), "", false)
	if err != nil {
		return errors.New("CreateMappingFileFailed")
	}

	unitConfFileContent, err := repo.webServerUnitFileFactory(
		vhostName,
		aliasesHostnames,
		publicDir,
		mappingFilePath,
	)
	if err != nil {
		return err
	}

	unitConfFilePathStr := infraEnvs.VirtualHostsConfDir + "/" + vhostFileNameStr
	unitConfFilePath, err := valueObject.NewUnixFilePath(unitConfFilePathStr)
	if err != nil {
		return errors.New(err.Error() + ": " + unitConfFilePathStr)
	}
	err = infraHelper.UpdateFile(
		unitConfFilePath.String(),
		unitConfFileContent,
		true,
	)
	if err != nil {
		return errors.New("CreateWebServerConfUnitFileFailed")
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

	aliases, err := repo.queryRepo.ReadAliasesByParentHostname(parentVhost.Hostname)
	if err != nil {
		return errors.New("GetParentVhostAliasesError: " + err.Error())
	}

	aliasesStr := []string{createDto.Hostname.String()}
	for _, alias := range aliases {
		aliasesStr = append(aliasesStr, alias.Hostname.String())
	}

	err = infraHelper.CreateSelfSignedSsl(
		infraEnvs.PkiConfDir,
		parentVhost.Hostname.String(),
		aliasesStr,
	)
	if err != nil {
		return errors.New("GenerateSelfSignedCertFailed")
	}

	err = repo.persistVirtualHost(createDto, parentVhost.RootDirectory)
	if err != nil {
		return err
	}

	return repo.createWebServerUnitFile(
		parentVhost.Hostname,
		parentVhost.RootDirectory,
	)
}

func (repo *VirtualHostCmdRepo) updateDirsOwnership(
	publicDir valueObject.UnixFilePath,
) error {
	directories := []string{
		publicDir.String(),
		"/app/conf/nginx",
		infraEnvs.PkiConfDir,
	}

	for _, directory := range directories {
		chownRecursively := true
		chownSymlinksToo := false
		err := infraHelper.UpdatePermissionsForWebServerUse(
			directory,
			chownRecursively,
			chownSymlinksToo,
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

	publicDirStr := infraEnvs.PrimaryPublicDir + "/" + hostnameStr
	publicDir, err := valueObject.NewUnixFilePath(publicDirStr)
	if err != nil {
		return errors.New(err.Error() + ": " + publicDirStr)
	}

	err = infraHelper.MakeDir(publicDirStr)
	if err != nil {
		return errors.New("MakePublicHtmlDirFailed")
	}

	aliases := []string{}
	err = infraHelper.CreateSelfSignedSsl(
		infraEnvs.PkiConfDir,
		hostnameStr,
		aliases,
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

	return repo.createWebServerUnitFile(createDto.Hostname, publicDir)
}

func (repo *VirtualHostCmdRepo) deleteWebServerUnitFile(
	vhostName valueObject.Fqdn,
) error {
	vhostFileNameStr := vhostName.String() + ".conf"
	if infraHelper.IsPrimaryVirtualHost(vhostName) {
		vhostFileNameStr = infraEnvs.PrimaryVhostFileName
	}

	mappingFilePathStr := infraEnvs.MappingsConfDir + "/" + vhostFileNameStr
	err := os.Remove(mappingFilePathStr)
	if err != nil {
		return err
	}

	webServerUnitFilePathStr := infraEnvs.VirtualHostsConfDir + "/" + vhostFileNameStr
	err = os.Remove(webServerUnitFilePathStr)
	if err != nil {
		return err
	}

	return infraHelper.ReloadWebServer()
}

func (repo *VirtualHostCmdRepo) Delete(vhost entity.VirtualHost) error {
	vhostNameStr := vhost.Hostname.String()
	err := repo.persistentDbSvc.Handler.
		Where(
			"hostname = ? OR parent_hostname = ?",
			vhostNameStr,
			vhostNameStr,
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

		return repo.createWebServerUnitFile(
			vhost.Hostname,
			vhost.RootDirectory,
		)
	}

	err = repo.deleteWebServerUnitFile(vhost.Hostname)
	if err != nil {
		return errors.New("DeleteWebServerUnitFileError: " + err.Error())
	}

	return nil
}
