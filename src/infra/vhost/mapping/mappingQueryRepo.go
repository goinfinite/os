package mappingInfra

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
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

func (repo *MappingQueryRepo) ReadById(
	id valueObject.MappingId,
) (entity entity.Mapping, err error) {
	query := dbModel.Mapping{
		ID: uint(id.Get()),
	}

	model := dbModel.Mapping{}
	err = repo.persistentDbSvc.Handler.
		Model(&query).
		First(&model).Error
	if err != nil {
		return entity, errors.New("ReadDatabaseEntryError")
	}

	entity, err = model.ToEntity()
	if err != nil {
		return entity, errors.New("ModelToEntityError")
	}

	return entity, nil
}

func (repo *MappingQueryRepo) ReadByHostname(
	hostname valueObject.Fqdn,
) (entities []entity.Mapping, err error) {
	query := dbModel.Mapping{
		Hostname: hostname.String(),
	}

	models := []dbModel.Mapping{}
	err = repo.persistentDbSvc.Handler.
		Model(query).
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

func (repo *MappingQueryRepo) ReadByServiceName(
	serviceName valueObject.ServiceName,
) (entities []entity.Mapping, err error) {
	svcNameStr := serviceName.String()
	query := dbModel.Mapping{
		TargetType:  "service",
		TargetValue: &svcNameStr,
	}

	models := []dbModel.Mapping{}
	err = repo.persistentDbSvc.Handler.
		Model(query).
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

func (repo *MappingQueryRepo) ReadWithMappings() (
	[]dto.VirtualHostWithMappings, error,
) {
	vhostsWithMappings := []dto.VirtualHostWithMappings{}

	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhosts, err := vhostQueryRepo.Get()
	if err != nil {
		return vhostsWithMappings, err
	}

	for _, vhost := range vhosts {
		mappings, err := repo.ReadByHostname(vhost.Hostname)
		if err != nil {
			log.Printf("[%s] ReadMappingsError: %s", vhost.Hostname, err.Error())
			continue
		}

		vhostsWithMappings = append(
			vhostsWithMappings,
			dto.NewVirtualHostWithMappings(
				vhost,
				mappings,
			),
		)
	}

	return vhostsWithMappings, nil
}
