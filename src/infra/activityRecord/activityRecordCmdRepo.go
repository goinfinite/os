package activityRecordInfra

import (
	"encoding/json"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
)

type ActivityRecordCmdRepo struct {
	trailDbSvc *internalDbInfra.TrailDatabaseService
}

func NewActivityRecordCmdRepo(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ActivityRecordCmdRepo {
	return &ActivityRecordCmdRepo{
		trailDbSvc: trailDbSvc,
	}
}

func (repo *ActivityRecordCmdRepo) Create(createDto dto.CreateActivityRecord) error {
	affectedResources := []dbModel.ActivityRecordAffectedResource{}
	for _, affectedResourceSri := range createDto.AffectedResources {
		affectedResourceModel := dbModel.ActivityRecordAffectedResource{
			SystemResourceIdentifier: affectedResourceSri.String(),
		}
		affectedResources = append(affectedResources, affectedResourceModel)
	}

	var recordDetails *string
	if createDto.RecordDetails != nil {
		recordDetailsBytes, err := json.Marshal(createDto.RecordDetails)
		if err != nil {
			return err
		}
		recordDetailsStr := string(recordDetailsBytes)
		recordDetails = &recordDetailsStr
	}

	var operatorAccountIdPtr *uint64
	if createDto.OperatorAccountId != nil {
		operatorAccountId := createDto.OperatorAccountId.Uint64()
		operatorAccountIdPtr = &operatorAccountId
	}

	var operatorIpAddressPtr *string
	if createDto.OperatorIpAddress != nil {
		operatorIpAddress := createDto.OperatorIpAddress.String()
		operatorIpAddressPtr = &operatorIpAddress
	}

	activityRecordModel := dbModel.NewActivityRecord(
		0, createDto.RecordLevel.String(), createDto.RecordCode.String(),
		affectedResources, recordDetails, operatorAccountIdPtr, operatorIpAddressPtr,
	)

	return repo.trailDbSvc.Handler.Create(&activityRecordModel).Error
}

func (repo *ActivityRecordCmdRepo) Delete(deleteDto dto.DeleteActivityRecord) error {
	deleteModel := dbModel.ActivityRecord{}
	if deleteDto.RecordId != nil {
		deleteModel.ID = deleteDto.RecordId.Uint64()
	}

	if deleteDto.RecordLevel != nil {
		deleteModel.RecordLevel = deleteDto.RecordLevel.String()
	}

	if deleteDto.RecordCode != nil {
		deleteModel.RecordCode = deleteDto.RecordCode.String()
	}

	affectedResources := []dbModel.ActivityRecordAffectedResource{}
	for _, affectedResourceSri := range deleteDto.AffectedResources {
		affectedResourceModel := dbModel.ActivityRecordAffectedResource{
			SystemResourceIdentifier: affectedResourceSri.String(),
		}
		affectedResources = append(affectedResources, affectedResourceModel)
	}
	deleteModel.AffectedResources = affectedResources

	if deleteDto.OperatorAccountId != nil {
		operatorAccountId := deleteDto.OperatorAccountId.Uint64()
		deleteModel.OperatorAccountId = &operatorAccountId
	}

	if deleteDto.OperatorIpAddress != nil {
		operatorIpAddressStr := deleteDto.OperatorIpAddress.String()
		deleteModel.OperatorIpAddress = &operatorIpAddressStr
	}

	dbQuery := repo.trailDbSvc.Handler.Model(&deleteModel).Where(&deleteModel)

	if deleteDto.CreatedBeforeAt != nil {
		dbQuery.Where("created_at < ?", deleteDto.CreatedBeforeAt.ReadAsGoTime())
	}
	if deleteDto.CreatedAfterAt != nil {
		dbQuery.Where("created_at > ?", deleteDto.CreatedAfterAt.ReadAsGoTime())
	}

	return dbQuery.Delete(&dbModel.ActivityRecord{}).Error
}
