package useCase

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func CreateSecurityActivityRecord(
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	code *valueObject.ActivityRecordCode,
	ipAddress *valueObject.IpAddress,
	operatorAccountId *valueObject.AccountId,
	targetAccountId *valueObject.AccountId,
	username *valueObject.Username,
) {
	recordLevel, _ := valueObject.NewActivityRecordLevel("SEC")

	createDto := dto.CreateActivityRecord{
		Level:             recordLevel,
		Code:              code,
		OperatorAccountId: operatorAccountId,
		TargetAccountId:   targetAccountId,
		IpAddress:         ipAddress,
		Username:          username,
	}

	go func() {
		_ = activityRecordCmdRepo.Create(createDto)
	}()
}
