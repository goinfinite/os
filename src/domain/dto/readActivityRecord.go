package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadActivityRecords struct {
	RecordId          *valueObject.ActivityRecordId          `json:"recordId,omitempty"`
	RecordLevel       *valueObject.ActivityRecordLevel       `json:"recordLevel,omitempty"`
	RecordCode        *valueObject.ActivityRecordCode        `json:"recordCode,omitempty"`
	AffectedResources []valueObject.SystemResourceIdentifier `json:"affectedResources,omitempty"`
	RecordDetails     *string                                `json:"recordDetails,omitempty"`
	OperatorAccountId *valueObject.AccountId                 `json:"operatorAccountId,omitempty"`
	OperatorIpAddress *valueObject.IpAddress                 `json:"operatorIpAddress,omitempty"`
	CreatedBeforeAt   *valueObject.UnixTime                  `json:"createdBeforeAt,omitempty"`
	CreatedAfterAt    *valueObject.UnixTime                  `json:"createdAfterAt,omitempty"`
}

func NewReadActivityRecords(
	recordId *valueObject.ActivityRecordId,
	recordLevel *valueObject.ActivityRecordLevel,
	recordCode *valueObject.ActivityRecordCode,
	affectedResources []valueObject.SystemResourceIdentifier,
	recordDetails *string,
	operatorAccountId *valueObject.AccountId,
	operatorIpAddress *valueObject.IpAddress,
	createdBeforeAt *valueObject.UnixTime,
	createdAfterAt *valueObject.UnixTime,
) ReadActivityRecords {
	return ReadActivityRecords{
		RecordId:          recordId,
		RecordLevel:       recordLevel,
		RecordCode:        recordCode,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
		AffectedResources: affectedResources,
		RecordDetails:     recordDetails,
		CreatedBeforeAt:   createdBeforeAt,
		CreatedAfterAt:    createdAfterAt,
	}
}
