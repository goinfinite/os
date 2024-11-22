package activityRecordInfra

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
)

type ActivityRecordQueryRepo struct {
	trailDbSvc *internalDbInfra.TrailDatabaseService
}

func NewActivityRecordQueryRepo(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ActivityRecordQueryRepo {
	return &ActivityRecordQueryRepo{
		trailDbSvc: trailDbSvc,
	}
}

func (repo *ActivityRecordQueryRepo) Read(
	readDto dto.ReadActivityRecords,
) ([]entity.ActivityRecord, error) {
	activityRecordEvents := []entity.ActivityRecord{}

	readModel := dbModel.ActivityRecord{}
	if readDto.RecordId != nil {
		recordId := readDto.RecordId.Uint64()
		readModel.ID = recordId
	}

	if readDto.RecordLevel != nil {
		recordLevelStr := readDto.RecordLevel.String()
		readModel.RecordLevel = recordLevelStr
	}

	if readDto.RecordCode != nil {
		recordCodeStr := readDto.RecordCode.String()
		readModel.RecordCode = recordCodeStr
	}

	if len(readDto.AffectedResources) > 0 {
		affectedResources := []dbModel.ActivityRecordAffectedResource{}
		for _, affectedResourceSri := range readDto.AffectedResources {
			affectedResourceModel := dbModel.ActivityRecordAffectedResource{
				SystemResourceIdentifier: affectedResourceSri.String(),
			}
			affectedResources = append(affectedResources, affectedResourceModel)
		}

		readModel.AffectedResources = affectedResources
	}

	if readDto.OperatorAccountId != nil {
		operatorAccountId := readDto.OperatorAccountId.Uint64()
		readModel.OperatorAccountId = &operatorAccountId
	}

	if readDto.OperatorIpAddress != nil {
		operatorIpAddressStr := readDto.OperatorIpAddress.String()
		readModel.OperatorIpAddress = &operatorIpAddressStr
	}

	dbQuery := repo.trailDbSvc.Handler.Model(&readModel).Where(&readModel)
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.GetAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.GetAsGoTime())
	}

	activityRecordEventModels := []dbModel.ActivityRecord{}
	err := dbQuery.
		Preload("AffectedResources").
		Find(&activityRecordEventModels).Error
	if err != nil {
		return activityRecordEvents, err
	}

	for _, activityRecordEventModel := range activityRecordEventModels {
		activityRecordEvent, err := activityRecordEventModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", activityRecordEventModel.ID),
				slog.Any("error", err),
			)
			continue
		}
		activityRecordEvents = append(activityRecordEvents, activityRecordEvent)
	}

	return activityRecordEvents, nil
}
