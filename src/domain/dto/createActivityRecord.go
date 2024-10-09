package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type CreateActivityRecord struct {
	Level             valueObject.ActivityRecordLevel    `json:"level"`
	Code              *valueObject.ActivityRecordCode    `json:"code,omitempty"`
	Message           *valueObject.ActivityRecordMessage `json:"message,omitempty"`
	IpAddress         *valueObject.IpAddress             `json:"ipAddress,omitempty"`
	OperatorAccountId *valueObject.AccountId             `json:"operatorAccountId,omitempty"`
	TargetAccountId   *valueObject.AccountId             `json:"targetAccountId,omitempty"`
	Username          *valueObject.Username              `json:"username,omitempty"`
	MappingId         *valueObject.MappingId             `json:"mappingId,omitempty"`
}
