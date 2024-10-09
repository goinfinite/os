package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type ActivityRecord struct {
	RecordId          valueObject.ActivityRecordId           `json:"recordId"`
	RecordLevel       valueObject.ActivityRecordLevel        `json:"recordLevel"`
	RecordCode        valueObject.ActivityRecordCode         `json:"recordCode,omitempty"`
	AffectedResources []valueObject.SystemResourceIdentifier `json:"affectedResources,omitempty"`
	RecordDetails     interface{}                            `json:"recordDetails,omitempty"`
	OperatorAccountId *valueObject.AccountId                 `json:"operatorAccountId,omitempty"`
	OperatorIpAddress *valueObject.IpAddress                 `json:"operatorIpAddress,omitempty"`
	CreatedAt         valueObject.UnixTime                   `json:"createdAt"`
}

func NewActivityRecord(
	recordId valueObject.ActivityRecordId,
	recordLevel valueObject.ActivityRecordLevel,
	recordCode valueObject.ActivityRecordCode,
	affectedResources []valueObject.SystemResourceIdentifier,
	recordDetails interface{},
	operatorAccountId *valueObject.AccountId,
	operatorIpAddress *valueObject.IpAddress,
	createdAt valueObject.UnixTime,
) (activityRecord ActivityRecord, err error) {
	return ActivityRecord{
		RecordId:          recordId,
		RecordLevel:       recordLevel,
		RecordCode:        recordCode,
		AffectedResources: affectedResources,
		RecordDetails:     recordDetails,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
		CreatedAt:         createdAt,
	}, nil
}
