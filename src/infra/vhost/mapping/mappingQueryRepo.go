package mappingInfra

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type MappingQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMappingQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MappingQueryRepo {
	return &MappingQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *MappingQueryRepo) GetById(
	id valueObject.MappingId,
) (entity.Mapping, error) {
	var mapping entity.Mapping

	mappingModel := dbModel.Mapping{ID: uint(id.Get())}
	err := repo.persistentDbSvc.Handler.Model(&mappingModel).
		First(&mappingModel).Error
	if err != nil {
		return mapping, errors.New("MappingNotFound")
	}

	return mappingModel.ToEntity()
}

func (repo *MappingQueryRepo) GetByHostname(
	hostname valueObject.Fqdn,
) ([]entity.Mapping, error) {
	mappingEntities := []entity.Mapping{}

	mappingModels := []dbModel.Mapping{}
	err := repo.persistentDbSvc.Handler.Model(
		dbModel.Mapping{Hostname: hostname.String()},
	).Find(&mappingModels).Error
	if err != nil {
		return mappingEntities, errors.New("DbQueryMappingsError")
	}

	for _, mappingModel := range mappingModels {
		mappingEntity, err := mappingModel.ToEntity()
		if err != nil {
			log.Printf("MappingModelToEntityError: %s", err.Error())
			continue
		}

		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return mappingEntities, nil
}
