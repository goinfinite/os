package activityRecordInfra

import (
	"errors"
	"log/slog"

	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkEntity "github.com/goinfinite/tk/src/domain/entity"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
	tkInfraDbModel "github.com/goinfinite/tk/src/infra/db/model"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"gorm.io/gorm"
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

func (repo *ActivityRecordQueryRepo) buildBaseQuery(
	readDto tkDto.ReadActivityRecordsRequest,
) *gorm.DB {
	readModel := tkInfraDbModel.ActivityRecord{}
	if readDto.RecordId != nil {
		readModel.ID = readDto.RecordId.Uint64()
	}

	if readDto.RecordLevel != nil {
		readModel.RecordLevel = readDto.RecordLevel.String()
	}

	if readDto.RecordCode != nil {
		readModel.RecordCode = readDto.RecordCode.String()
	}

	if len(readDto.AffectedResources) > 0 {
		affectedResources := []tkInfraDbModel.ActivityRecordAffectedResource{}
		for _, affectedResourceSri := range readDto.AffectedResources {
			affectedResourceModel := tkInfraDbModel.ActivityRecordAffectedResource{
				SystemResourceIdentifier: affectedResourceSri.String(),
			}
			affectedResources = append(affectedResources, affectedResourceModel)
		}
		readModel.AffectedResources = affectedResources
	}

	if readDto.OperatorSri != nil {
		operatorSriStr := readDto.OperatorSri.String()
		readModel.OperatorSri = &operatorSriStr
	}

	if readDto.OperatorIpAddress != nil {
		operatorIpAddressStr := readDto.OperatorIpAddress.String()
		readModel.OperatorIpAddress = &operatorIpAddressStr
	}

	dbQuery := repo.trailDbSvc.Handler.Model(&readModel).Where(&readModel)
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.ReadAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.ReadAsGoTime())
	}

	return dbQuery
}

func (repo *ActivityRecordQueryRepo) Read(
	readDto tkDto.ReadActivityRecordsRequest,
) (responseDto tkDto.ReadActivityRecordsResponse, err error) {
	dbQuery := repo.buildBaseQuery(readDto)

	dbQuery, responsePagination, err := tkInfraDb.PaginationQueryBuilder(
		dbQuery, readDto.Pagination,
	)
	if err != nil {
		return responseDto, err
	}
	responseDto.Pagination = responsePagination

	activityRecordModels := []tkInfraDbModel.ActivityRecord{}
	err = dbQuery.
		Preload("AffectedResources").
		Find(&activityRecordModels).Error
	if err != nil {
		return responseDto, err
	}

	for _, activityRecordModel := range activityRecordModels {
		activityRecordEntity, err := activityRecordModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", activityRecordModel.ID),
				slog.String("err", err.Error()),
			)
			continue
		}
		responseDto.ActivityRecords = append(
			responseDto.ActivityRecords, activityRecordEntity,
		)
	}

	return responseDto, nil
}

func (repo *ActivityRecordQueryRepo) ReadFirst(
	readDto tkDto.ReadActivityRecordsRequest,
) (tkEntity.ActivityRecord, error) {
	dbQuery := repo.buildBaseQuery(readDto)

	activityRecordModel := tkInfraDbModel.ActivityRecord{}
	err := dbQuery.
		Preload("AffectedResources").
		First(&activityRecordModel).Error
	if err != nil {
		return tkEntity.ActivityRecord{}, errors.New("ActivityRecordNotFound")
	}

	return activityRecordModel.ToEntity()
}
