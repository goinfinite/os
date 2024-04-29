package mappingInfra

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
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

func (repo *MappingCmdRepo) mappingConfigFactory(
	mapping entity.Mapping,
) (string, error) {
	mappingConfig := "location " + mapping.Path.String() + " {"

	switch mapping.TargetType.String() {
	case "url":
		mappingConfig += `
	return 301 ` + mapping.TargetUrl.String() + `;`
	case "service":
		mappingConfig += ``
	case "response-code":
		mappingConfig += `
	return ` + mapping.TargetHttpResponseCode.String() + `;`
	case "inline-html":
		mappingConfig += `
	add_header Content-Type text/html;
	return 200 ` + mapping.TargetInlineHtmlContent.String() + `;`
	case "static-files":
		mappingConfig += `
	try_files $uri $uri/ index.html?$query_string;`
	}

	mappingConfig += `
}
`
	return mappingConfig, nil
}

func (repo *MappingCmdRepo) rebuildMappingFile(
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

	fullMappingConfigContent := ""
	for _, mapping := range mappings {
		mappingConfigContent, err := repo.mappingConfigFactory(mapping)
		if err != nil {
			log.Printf(
				"MappingConfigFactoryError (%s): %s",
				mapping.Path.String(),
				err.Error(),
			)
		}
		fullMappingConfigContent += mappingConfigContent
	}

	shouldOverwrite := true
	return infraHelper.UpdateFile(
		mappingFilePath.String(),
		fullMappingConfigContent,
		shouldOverwrite,
	)
}

func (repo *MappingCmdRepo) Create(
	createDto dto.CreateMapping,
) (valueObject.MappingId, error) {
	var mappingId valueObject.MappingId

	isServiceMapping := createDto.TargetType.String() == "service"
	isPhpServiceMapping := isServiceMapping && createDto.TargetServiceName.String() == "php"
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

	err = repo.rebuildMappingFile(createDto.Hostname)
	if err != nil {
		return mappingId, err
	}

	return mappingId, repo.vhostCmdRepo.ReloadWebServer()
}

func (repo *MappingCmdRepo) DeleteMapping(mappingId valueObject.MappingId) error {
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

	err = repo.rebuildMappingFile(mapping.Hostname)
	if err != nil {
		return err
	}

	return repo.vhostCmdRepo.ReloadWebServer()
}

func (repo *MappingCmdRepo) DeleteAutoMapping(
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

		targetServiceName := primaryVhostMapping.TargetServiceName
		if targetServiceName == nil {
			continue
		}

		if targetServiceName.String() != serviceName.String() {
			continue
		}

		mappingIdToDelete = &primaryVhostMapping.Id
	}

	if mappingIdToDelete == nil {
		return nil
	}

	return repo.DeleteMapping(*mappingIdToDelete)
}
