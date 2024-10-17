package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type DeleteActivityRecord struct {
	RecordId          *valueObject.ActivityRecordId          `json:"recordId,omitempty"`
	RecordLevel       *valueObject.ActivityRecordLevel       `json:"recordLevel,omitempty"`
	RecordCode        *valueObject.ActivityRecordCode        `json:"recordCode,omitempty"`
	AffectedResources []valueObject.SystemResourceIdentifier `json:"affectedResources,omitempty"`
	OperatorAccountId *valueObject.AccountId                 `json:"operatorAccountId,omitempty"`
	OperatorIpAddress *valueObject.IpAddress                 `json:"operatorIpAddress,omitempty"`
	CreatedBeforeAt   *valueObject.UnixTime                  `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt    *valueObject.UnixTime                  `json:"createdAfterAt,omitempty"`
}

func NewDeleteActivityRecord(
	recordId *valueObject.ActivityRecordId,
	recordLevel *valueObject.ActivityRecordLevel,
	recordCode *valueObject.ActivityRecordCode,
	affectedResources []valueObject.SystemResourceIdentifier,
	operatorAccountId *valueObject.AccountId,
	operatorIpAddress *valueObject.IpAddress,
	createdBeforeAt, createdAfterAt *valueObject.UnixTime,
) DeleteActivityRecord {
	return DeleteActivityRecord{
		RecordId:          recordId,
		RecordLevel:       recordLevel,
		RecordCode:        recordCode,
		AffectedResources: affectedResources,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
		CreatedBeforeAt:   createdBeforeAt,
		CreatedAfterAt:    createdAfterAt,
	}
}
