package vhostInfra

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type VirtualHostQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewVirtualHostQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostQueryRepo {
	return &VirtualHostQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *VirtualHostQueryRepo) Read() ([]entity.VirtualHost, error) {
	entities := []entity.VirtualHost{}

	models := []dbModel.VirtualHost{}
	err := repo.persistentDbSvc.Handler.
		Model(&models).
		Find(&models).Error
	if err != nil {
		return entities, errors.New("ReadDatabaseEntriesError")
	}

	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			log.Printf("ModelToEntityError: %s", err.Error())
			continue
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (repo *VirtualHostQueryRepo) ReadByHostname(
	hostname valueObject.Fqdn,
) (vhostEntity entity.VirtualHost, err error) {

	model := dbModel.VirtualHost{}
	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.VirtualHost{}).
		Where("hostname = ?", hostname.String()).
		First(&model).Error
	if err != nil {
		errorMessage := "VhostNotFound"
		if err.Error() != "record not found" {
			errorMessage = "ReadDatabaseEntryError"
		}

		return vhostEntity, errors.New(errorMessage)
	}

	vhostEntity, err = model.ToEntity()
	if err != nil {
		return vhostEntity, errors.New("ModelToEntityError")
	}

	return vhostEntity, nil
}

func (repo *VirtualHostQueryRepo) ReadAliasesByParentHostname(
	hostname valueObject.Fqdn,
) ([]entity.VirtualHost, error) {
	aliasesEntities := []entity.VirtualHost{}

	aliasesModels := []dbModel.VirtualHost{}
	err := repo.persistentDbSvc.Handler.
		Model(&aliasesModels).
		Where("parent_hostname = ?", hostname.String()).
		Find(&aliasesModels).Error
	if err != nil {
		return aliasesEntities, errors.New("ReadDatabaseEntriesError")
	}

	for _, aliasModel := range aliasesModels {
		aliasEntity, err := aliasModel.ToEntity()
		if err != nil {
			log.Printf("ModelToEntityError: %s", err.Error())
			continue
		}

		aliasesEntities = append(aliasesEntities, aliasEntity)
	}

	return aliasesEntities, nil
}

func (repo *VirtualHostQueryRepo) GetVirtualHostMappingsFilePath(
	vhostName valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	var vhostFilePath valueObject.UnixFilePath

	vhost, err := repo.ReadByHostname(vhostName)
	if err != nil {
		return vhostFilePath, errors.New("VirtualHostNotFound")
	}

	if vhost.Type.String() == "alias" {
		vhostName = *vhost.ParentHostname
	}

	vhostFileNameStr := vhostName.String() + ".conf"
	if infraHelper.IsPrimaryVirtualHost(vhostName) {
		vhostFileNameStr = infraEnvs.PrimaryVhostFileName
	}

	return valueObject.NewUnixFilePath(
		infraEnvs.MappingsConfDir + "/" + vhostFileNameStr,
	)
}
