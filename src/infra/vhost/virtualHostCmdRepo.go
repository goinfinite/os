package vhostInfra

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/goinfinite/os/src/domain/dto"
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
	vhostHostname valueObject.Fqdn,
	aliasesHostnames []valueObject.Fqdn,
	publicDir valueObject.UnixFilePath,
	mappingFilePath valueObject.UnixFilePath,
) (string, error) {
	vhostHostnameStr := vhostHostname.String()

	aliasesHostnamesStr := []string{}
	for _, aliasHostname := range aliasesHostnames {
		aliasesHostnamesStr = append(aliasesHostnamesStr, aliasHostname.String())
	}

	valuesToInterpolate := map[string]interface{}{
		"VirtualHostHostname": vhostHostnameStr,
		"AliasesHostnames":    aliasesHostnamesStr,
		"PublicDirectory":     publicDir,
		"CertPath":            infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".crt",
		"KeyPath":             infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".key",
		"MappingFilePath":     mappingFilePath,
	}

	unitConfTemplate := `server {
    listen 80;
    listen 443 ssl;
    server_name {{ .VirtualHostHostname }} www.{{ .VirtualHostHostname }}{{ range $aliasHostname := .AliasesHostnames }} {{ $aliasHostname }} www.{{ $aliasHostname }}{{ end }};

    root {{ .PublicDirectory }};

    ssl_certificate {{ .CertPath }};
    ssl_certificate_key {{ .KeyPath }};

    access_log /app/logs/nginx/{{ .VirtualHostHostname }}_access.log combined buffer=512k flush=1m;
    error_log /app/logs/nginx/{{ .VirtualHostHostname }}_error.log warn;

    include /etc/nginx/std.conf;
    include {{ .MappingFilePath }};
}
`

	unitConfTemplatePtr, err := template.New("webServerConfUnitFile").Parse(unitConfTemplate)
	if err != nil {
		return "", errors.New("TemplateParsingError: " + err.Error())
	}

	var unitConfFileContent strings.Builder
	err = unitConfTemplatePtr.Execute(&unitConfFileContent, valuesToInterpolate)
	if err != nil {
		return "", errors.New("TemplateExecutionError: " + err.Error())
	}

	return unitConfFileContent.String(), nil
}

func (repo *VirtualHostCmdRepo) createWebServerUnitFile(
	vhostHostname valueObject.Fqdn,
) error {
	vhostEntity, err := repo.queryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &vhostHostname,
	})
	if err != nil {
		return errors.New("ReadVirtualHostEntityError: " + err.Error())
	}

	if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
		if vhostEntity.ParentHostname == nil {
			return errors.New("AliasMissingParentHostname")
		}

		vhostEntity, err = repo.queryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
			Hostname: vhostEntity.ParentHostname,
		})
		if err != nil {
			return errors.New("ReadAliasParentVirtualHostError: " + err.Error())
		}
		vhostHostname = vhostEntity.Hostname
	}

	mappingsFilePath, err := repo.queryRepo.ReadVirtualHostMappingsFilePath(vhostHostname)
	if err != nil {
		return errors.New("ReadVirtualHostMappingsFilePathError: " + err.Error())
	}

	err = infraHelper.UpdateFile(mappingsFilePath.String(), "", false)
	if err != nil {
		return errors.New("TruncateMappingFileFailed: " + err.Error())
	}

	unitConfFileContent, err := repo.webServerUnitFileFactory(
		vhostEntity.Hostname, vhostEntity.AliasesHostnames,
		vhostEntity.RootDirectory, mappingsFilePath,
	)
	if err != nil {
		return err
	}

	mappingsFileNameStr := mappingsFilePath.ReadFileName().String()
	rawUnitConfFilePath := infraEnvs.VirtualHostsConfDir + "/" + mappingsFileNameStr
	unitConfFilePath, err := valueObject.NewUnixFilePath(rawUnitConfFilePath)
	if err != nil {
		return errors.New("InvalidUnitConfFilePath")
	}

	err = infraHelper.UpdateFile(unitConfFilePath.String(), unitConfFileContent, true)
	if err != nil {
		return errors.New("CreateWebServerConfUnitFileFailed: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}

func (repo *VirtualHostCmdRepo) createVirtualHostPublicDirectory(
	createDto dto.CreateVirtualHost,
) (publicDir valueObject.UnixFilePath, err error) {
	if createDto.Type == valueObject.VirtualHostTypeAlias {
		parentVirtualHostEntity, err := repo.queryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
			Hostname: createDto.ParentHostname,
		})
		if err != nil {
			return publicDir, errors.New("ReadAliasParentVirtualHostError: " + err.Error())
		}

		return parentVirtualHostEntity.RootDirectory, nil
	}

	rawPublicDir := infraEnvs.PrimaryPublicDir + "/" + createDto.Hostname.String()

	publicDir, err = valueObject.NewUnixFilePath(rawPublicDir)
	if err != nil {
		return publicDir, errors.New("InvalidVirtualHostPublicDir")
	}

	err = infraHelper.MakeDir(publicDir.String())
	if err != nil {
		return publicDir, errors.New("CreateVirtualHostPublicDirFailed")
	}

	return publicDir, nil
}

func (repo *VirtualHostCmdRepo) Create(createDto dto.CreateVirtualHost) error {
	publicDir, err := repo.createVirtualHostPublicDirectory(createDto)
	if err != nil {
		return errors.New("CreateVirtualHostPublicDirFailed: " + err.Error())
	}

	pkiConfDir, err := valueObject.NewUnixFilePath(infraEnvs.PkiConfDir)
	if err != nil {
		return errors.New("InvalidPkiConfDir")
	}

	if createDto.Type != valueObject.VirtualHostTypeAlias {
		aliasHostnames := []valueObject.Fqdn{}
		err = infraHelper.CreateSelfSignedSsl(pkiConfDir, createDto.Hostname, aliasHostnames)
		if err != nil {
			return errors.New("CreateSelfSignedSslFailed: " + err.Error())
		}
	}

	webServerConfDir, err := valueObject.NewUnixFilePath(infraEnvs.VirtualHostsConfDir)
	if err != nil {
		return errors.New("InvalidWebServerConfDir")
	}

	vhostRelatedDirectories := []valueObject.UnixFilePath{
		publicDir, pkiConfDir, webServerConfDir,
	}
	for _, directory := range vhostRelatedDirectories {
		chownRecursively := true
		chownSymlinksToo := false
		err := infraHelper.UpdateOwnershipForWebServerUse(
			directory.String(), chownRecursively, chownSymlinksToo,
		)
		if err != nil {
			return errors.New("UpdateOwnershipForWebServerUseError: " + err.Error())
		}
	}

	virtualHostModel := dbModel.VirtualHost{
		Hostname:      createDto.Hostname.String(),
		Type:          createDto.Type.String(),
		RootDirectory: publicDir.String(),
	}
	if createDto.ParentHostname != nil {
		parentHostnameStr := createDto.ParentHostname.String()
		virtualHostModel.ParentHostname = &parentHostnameStr
	}

	err = repo.persistentDbSvc.Handler.Create(&virtualHostModel).Error
	if err != nil {
		return errors.New("DbCreateVirtualHostError: " + err.Error())
	}

	return repo.createWebServerUnitFile(createDto.Hostname)
}

func (repo *VirtualHostCmdRepo) deleteWebServerUnitFile(
	vhostHostname valueObject.Fqdn,
) error {
	mappingsFilePath, err := repo.queryRepo.ReadVirtualHostMappingsFilePath(vhostHostname)
	if err != nil {
		return errors.New("ReadVirtualHostMappingsFilePathError: " + err.Error())
	}

	err = os.Remove(mappingsFilePath.String())
	if err != nil {
		return err
	}

	mappingsFileNameStr := mappingsFilePath.ReadFileName().String()
	webServerUnitFilePathStr := infraEnvs.VirtualHostsConfDir + "/" + mappingsFileNameStr
	err = os.Remove(webServerUnitFilePathStr)
	if err != nil {
		return err
	}

	return infraHelper.ReloadWebServer()
}

func (repo *VirtualHostCmdRepo) Delete(vhostHostname valueObject.Fqdn) error {
	vhostEntity, err := repo.queryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &vhostHostname,
	})
	if err != nil {
		return errors.New("ReadVirtualHostEntityError: " + err.Error())
	}

	vhostHostnameStr := vhostHostname.String()
	err = repo.persistentDbSvc.Handler.
		Where("hostname = ? OR parent_hostname = ?", vhostHostnameStr, vhostHostnameStr).
		Delete(dbModel.VirtualHost{}).Error
	if err != nil {
		return err
	}

	if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
		if vhostEntity.ParentHostname == nil {
			return errors.New("AliasMissingParentHostname")
		}

		parentVirtualHostEntity, err := repo.queryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
			Hostname: vhostEntity.ParentHostname,
		})
		if err != nil {
			return errors.New("ReadAliasParentVirtualHostError: " + err.Error())
		}

		return repo.createWebServerUnitFile(parentVirtualHostEntity.Hostname)
	}

	pkiConfDir, err := valueObject.NewUnixFilePath(infraEnvs.PkiConfDir)
	if err != nil {
		return errors.New("InvalidPkiConfDir")
	}
	pkiConfDirStr := pkiConfDir.String()

	vhostCertFilePath := pkiConfDirStr + "/" + vhostHostnameStr + ".crt"
	err = os.Remove(vhostCertFilePath)
	if err != nil {
		slog.Error("RemoveSslCertFileError", slog.String("error", err.Error()))
	}

	vhostCertKeyFilePath := pkiConfDirStr + "/" + vhostHostnameStr + ".key"
	err = os.Remove(vhostCertKeyFilePath)
	if err != nil {
		slog.Error("RemoveSslCertKeyFileError", slog.String("error", err.Error()))
	}

	return repo.deleteWebServerUnitFile(vhostHostname)
}
