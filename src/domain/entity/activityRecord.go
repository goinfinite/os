package entity

import (
	"errors"

	"github.com/speedianet/os/src/domain/valueObject"
)

type ActivityRecord struct {
	Id                valueObject.ActivityRecordId       `json:"id"`
	Level             valueObject.ActivityRecordLevel    `json:"level"`
	Code              *valueObject.ActivityRecordCode    `json:"code,omitempty"`
	Message           *valueObject.ActivityRecordMessage `json:"message,omitempty"`
	IpAddress         *valueObject.IpAddress             `json:"ipAddress,omitempty"`
	OperatorAccountId *valueObject.AccountId             `json:"operatorAccountId,omitempty"`
	TargetAccountId   *valueObject.AccountId             `json:"targetAccountId,omitempty"`
	Username          *valueObject.Username              `json:"username,omitempty"`
	MappingId         *valueObject.MappingId             `json:"mappingId,omitempty"`
	CreatedAt         valueObject.UnixTime               `json:"createdAt"`
}

func NewActivityRecord(
	id valueObject.ActivityRecordId,
	level valueObject.ActivityRecordLevel,
	code *valueObject.ActivityRecordCode,
	message *valueObject.ActivityRecordMessage,
	ipAddress *valueObject.IpAddress,
	operatorAccountId *valueObject.AccountId,
	targetAccountId *valueObject.AccountId,
	username *valueObject.Username,
	mappingId *valueObject.MappingId,
	createdAt valueObject.UnixTime,
) (activityRecord ActivityRecord, err error) {
	if code == nil && message == nil {
		return activityRecord, errors.New("CodeOrMessageRequired")
	}

	return ActivityRecord{
		Id:                id,
		Level:             level,
		Code:              code,
		Message:           message,
		IpAddress:         ipAddress,
		OperatorAccountId: operatorAccountId,
		TargetAccountId:   targetAccountId,
		CreatedAt:         createdAt,
	}, nil
}
