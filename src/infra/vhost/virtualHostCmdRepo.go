package vhostInfra

import (
	"errors"
	"log/slog"
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
	vhostQueryRepo  *VirtualHostQueryRepo
}

func NewVirtualHostCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostCmdRepo {
	return &VirtualHostCmdRepo{
		persistentDbSvc: persistentDbSvc,
		vhostQueryRepo:  NewVirtualHostQueryRepo(persistentDbSvc),
	}
}

func (repo *VirtualHostCmdRepo) webServerUnitFileFactory(
	vhostEntity entity.VirtualHost,
	mappingFilePath valueObject.UnixFilePath,
) (string, error) {
	vhostHostnameStr := vhostEntity.Hostname.String()

	aliasesHostnamesStr := []string{}
	for _, aliasHostname := range vhostEntity.AliasesHostnames {
		aliasesHostnamesStr = append(aliasesHostnamesStr, aliasHostname.String())
	}

	mainServerName := vhostHostnameStr + " www." + vhostHostnameStr
	if vhostEntity.IsWildcard || vhostEntity.Type == valueObject.VirtualHostTypeWildcard {
		mainServerName += " *." + vhostHostnameStr
	}

	confVariables := map[string]interface{}{
		"VirtualHostHostname": vhostHostnameStr,
		"MainServerName":      mainServerName,
		"AliasesHostnames":    aliasesHostnamesStr,
		"PublicDirectory":     vhostEntity.RootDirectory.String(),
		"CertPath":            infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".crt",
		"KeyPath":             infraEnvs.PkiConfDir + "/" + vhostHostnameStr + ".key",
		"MappingFilePath":     mappingFilePath.String(),
	}

	unitConfTemplate := `server {
    listen 80;
    listen 443 ssl;
    server_name {{ .MainServerName }}{{ range $aliasHostname := .AliasesHostnames }} {{ $aliasHostname }} www.{{ $aliasHostname }}{{ end }};

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
	err = unitConfTemplatePtr.Execute(&unitConfFileContent, confVariables)
	if err != nil {
		return "", errors.New("TemplateExecutionError: " + err.Error())
	}

	return unitConfFileContent.String(), nil
}

func (repo *VirtualHostCmdRepo) ReadVirtualHostWebServerUnitFileFilePath(
	vhostHostname valueObject.Fqdn,
) (unitFilePath valueObject.UnixFilePath, err error) {
	mappingsFilePath, err := repo.vhostQueryRepo.ReadVirtualHostMappingsFilePath(vhostHostname)
	if err != nil {
		return unitFilePath, errors.New("ReadVirtualHostMappingsFilePathError: " + err.Error())
	}

	mappingsFileNameStr := mappingsFilePath.ReadFileName().String()
	rawUnitConfFilePath := infraEnvs.VirtualHostsConfDir + "/" + mappingsFileNameStr
	return valueObject.NewUnixFilePath(rawUnitConfFilePath)
}

func (repo *VirtualHostCmdRepo) createWebServerUnitFile(
	vhostHostname valueObject.Fqdn,
) error {
	vhostEntity, err := repo.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &vhostHostname,
	})
	if err != nil {
		return errors.New("ReadVirtualHostEntityError: " + err.Error())
	}

	if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
		if vhostEntity.ParentHostname == nil {
			return errors.New("AliasMissingParentHostname")
		}

		vhostEntity, err = repo.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
			Hostname: vhostEntity.ParentHostname,
		})
		if err != nil {
			return errors.New("ReadAliasParentVirtualHostError: " + err.Error())
		}
		vhostHostname = vhostEntity.Hostname
	}

	mappingsFilePath, err := repo.vhostQueryRepo.ReadVirtualHostMappingsFilePath(vhostHostname)
	if err != nil {
		return errors.New("ReadVirtualHostMappingsFilePathError: " + err.Error())
	}

	unitConfFileContent, err := repo.webServerUnitFileFactory(
		vhostEntity, mappingsFilePath,
	)
	if err != nil {
		return err
	}

	unitConfFilePath, err := repo.ReadVirtualHostWebServerUnitFileFilePath(vhostHostname)
	if err != nil {
		return errors.New("ReadWebServerUnitConfFilePathError: " + err.Error())
	}

	err = infraHelper.UpdateFile(unitConfFilePath.String(), unitConfFileContent, true)
	if err != nil {
		return errors.New("CreateWebServerConfUnitFileFailed: " + err.Error())
	}

	mappingCmdRepo := NewMappingCmdRepo(repo.persistentDbSvc)
	err = mappingCmdRepo.RecreateMappingFile(vhostHostname)
	if err != nil {
		return errors.New("RecreateMappingFileError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}

func (repo *VirtualHostCmdRepo) createVirtualHostPublicDirectory(
	createDto dto.CreateVirtualHost,
) (publicDir valueObject.UnixFilePath, err error) {
	if createDto.Type == valueObject.VirtualHostTypeAlias {
		parentVirtualHostEntity, err := repo.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
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

	isWildcard := false
	if createDto.IsWildcard != nil {
		isWildcard = *createDto.IsWildcard
	}
	virtualHostModel := dbModel.VirtualHost{
		Hostname:      createDto.Hostname.String(),
		Type:          createDto.Type.String(),
		RootDirectory: publicDir.String(),
		IsPrimary:     false,
		IsWildcard:    isWildcard,
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

func (repo *VirtualHostCmdRepo) Delete(vhostHostname valueObject.Fqdn) error {
	withMappings := true
	vhostReadResponse, err := repo.vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination:   dto.PaginationSingleItem,
		Hostname:     &vhostHostname,
		WithMappings: &withMappings,
	})
	if err != nil {
		return errors.New("ReadVirtualHostEntityError: " + err.Error())
	}

	if len(vhostReadResponse.VirtualHostWithMappings) == 0 {
		return errors.New("VirtualHostNotFound")
	}

	mappingCmdRepo := NewMappingCmdRepo(repo.persistentDbSvc)
	for _, mappingEntity := range vhostReadResponse.VirtualHostWithMappings[0].Mappings {
		err = mappingCmdRepo.Delete(mappingEntity.Id)
		if err != nil {
			slog.Error(
				"DeleteMappingError",
				slog.String("mappingId", mappingEntity.Id.String()),
				slog.String("error", err.Error()),
			)
		}
	}

	vhostWebServerConfFilePath, err := repo.ReadVirtualHostWebServerUnitFileFilePath(
		vhostHostname,
	)
	if err != nil {
		return errors.New("ReadWebServerUnitConfFilePathError: " + err.Error())
	}

	vhostHostnameStr := vhostHostname.String()
	err = repo.persistentDbSvc.Handler.
		Where("hostname = ? OR parent_hostname = ?", vhostHostnameStr, vhostHostnameStr).
		Delete(dbModel.VirtualHost{}).Error
	if err != nil {
		return err
	}

	vhostEntity := vhostReadResponse.VirtualHostWithMappings[0].VirtualHost
	if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
		if vhostEntity.ParentHostname == nil {
			return errors.New("AliasMissingParentHostname")
		}

		parentVirtualHostEntity, err := repo.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
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

	err = os.Remove(vhostWebServerConfFilePath.String())
	if err != nil {
		return errors.New("RemoveWebServerUnitConfFileError: " + err.Error())
	}

	return infraHelper.ReloadWebServer()
}
