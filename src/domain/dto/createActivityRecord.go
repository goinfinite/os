package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type CreateActivityRecord struct {
	RecordLevel       valueObject.ActivityRecordLevel        `json:"recordLevel"`
	RecordCode        valueObject.ActivityRecordCode         `json:"recordCode"`
	AffectedResources []valueObject.SystemResourceIdentifier `json:"affectedResources,omitempty"`
	RecordDetails     interface{}                            `json:"recordDetails,omitempty"`
	OperatorAccountId *valueObject.AccountId                 `json:"operatorAccountId,omitempty"`
	OperatorIpAddress *valueObject.IpAddress                 `json:"operatorIpAddress,omitempty"`
}
