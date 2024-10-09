package dbModel

import (
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ActivityRecord struct {
	ID                uint64 `gorm:"primarykey"`
	RecordLevel       string `gorm:"not null"`
	RecordCode        string `gorm:"not null"`
	AffectedResources []ActivityRecordAffectedResource
	RecordDetails     *string
	OperatorAccountId *uint64
	OperatorIpAddress *string
	CreatedAt         time.Time `gorm:"not null"`
}

func (ActivityRecord) TableName() string {
	return "activity_records"
}

func NewActivityRecord(
	recordId uint64,
	recordLevel, recordCode string,
	affectedResources []ActivityRecordAffectedResource,
	recordDetails *string,
	operatorAccountId *uint64,
	operatorIpAddress *string,
) ActivityRecord {
	model := ActivityRecord{
		RecordLevel:       recordLevel,
		RecordCode:        recordCode,
		AffectedResources: affectedResources,
		RecordDetails:     recordDetails,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}

	if recordId != 0 {
		model.ID = recordId
	}

	return model
}

func (model ActivityRecord) ToEntity() (recordEntity entity.ActivityRecord, err error) {
	recordId, err := valueObject.NewActivityRecordId(model.ID)
	if err != nil {
		return recordEntity, err
	}

	recordLevel, err := valueObject.NewActivityRecordLevel(model.RecordLevel)
	if err != nil {
		return recordEntity, err
	}

	recordCode, err := valueObject.NewActivityRecordCode(model.RecordCode)
	if err != nil {
		return recordEntity, err
	}

	affectedResources := []valueObject.SystemResourceIdentifier{}
	for _, resource := range model.AffectedResources {
		sri, err := valueObject.NewSystemResourceIdentifier(resource.SystemResourceIdentifier)
		if err != nil {
			return recordEntity, err
		}
		affectedResources = append(affectedResources, sri)
	}

	var recordDetails interface{}
	if model.RecordDetails != nil {
		recordDetails = *model.RecordDetails
	}

	var operatorAccountIdPtr *valueObject.AccountId
	if model.OperatorAccountId != nil {
		operatorAccountId, err := valueObject.NewAccountId(*model.OperatorAccountId)
		if err != nil {
			return recordEntity, err
		}
		operatorAccountIdPtr = &operatorAccountId
	}

	var operatorIpAddressPtr *valueObject.IpAddress
	if model.OperatorIpAddress != nil {
		operatorIpAddress, err := valueObject.NewIpAddress(*model.OperatorIpAddress)
		if err != nil {
			return recordEntity, err
		}
		operatorIpAddressPtr = &operatorIpAddress
	}

	createdAt := valueObject.NewUnixTimeWithGoTime(model.CreatedAt)

	return entity.NewActivityRecord(
		recordId, recordLevel, recordCode, affectedResources, recordDetails,
		operatorAccountIdPtr, operatorIpAddressPtr, createdAt,
	)
}
