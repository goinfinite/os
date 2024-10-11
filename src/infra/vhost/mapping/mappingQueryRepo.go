package mappingInfra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
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
	model := dbModel.Mapping{}
	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.Mapping{}).
		Where("id = ?", id.Uint64()).
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
) ([]entity.Mapping, error) {
	entities := []entity.Mapping{}

	models := []dbModel.Mapping{}
	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.Mapping{}).
		Where("hostname = ?", hostname.String()).
		Find(&models).Error
	if err != nil {
		return entities, errors.New("ReadDatabaseEntriesError")
	}

	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			slog.Error("ModelToEntityError", slog.Any("error", err))
			continue
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (repo *MappingQueryRepo) ReadByServiceName(
	serviceName valueObject.ServiceName,
) (entities []entity.Mapping, err error) {
	models := []dbModel.Mapping{}
	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.Mapping{}).
		Where("target_type = 'service' AND target_value = ?", serviceName.String()).
		Find(&models).Error
	if err != nil {
		return entities, errors.New("ReadDatabaseEntriesError")
	}

	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			slog.Error("ModelToEntityError", slog.Any("error", err))
			continue
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (repo *MappingQueryRepo) ReadWithMappings() (
	vhostsWithMappings []dto.VirtualHostWithMappings, err error,
) {
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(repo.persistentDbSvc)
	vhosts, err := vhostQueryRepo.Read()
	if err != nil {
		return vhostsWithMappings, err
	}

	for _, vhost := range vhosts {
		mappings, err := repo.ReadByHostname(vhost.Hostname)
		if err != nil {
			slog.Error(
				"ReadMappingsError",
				slog.String("vhostHostname", vhost.Hostname.String()),
				slog.Any("error", err),
			)
			continue
		}

		vhostsWithMappings = append(
			vhostsWithMappings,
			dto.NewVirtualHostWithMappings(vhost, mappings),
		)
	}

	return vhostsWithMappings, nil
}
