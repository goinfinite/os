package useCase

import (
	"errors"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

const MaxFailedLoginAttemptsPerIpAddress uint = 3
const FailedLoginAttemptsInterval time.Duration = 15 * time.Minute
const SessionTokenExpiresIn time.Duration = 3 * time.Hour

func getFailedLoginAttemptsCount(
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	createDto dto.CreateSessionToken,
) uint {
	recordLevel, _ := valueObject.NewActivityRecordLevel("SEC")
	recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
	failedAttemptsIntervalStartsAt := valueObject.NewUnixTimeBeforeNow(
		FailedLoginAttemptsInterval,
	)
	readDto := dto.NewReadActivityRecords(
		nil, &recordLevel, &recordCode, nil, nil, nil, &createDto.OperatorIpAddress,
		nil, &failedAttemptsIntervalStartsAt,
	)

	failedLoginAttempts := ReadActivityRecords(activityRecordQueryRepo, readDto)
	return uint(len(failedLoginAttempts))
}

func CreateSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateSessionToken,
) (accessToken entity.AccessToken, err error) {
	failedAttemptsCount := getFailedLoginAttemptsCount(activityRecordQueryRepo, createDto)
	if failedAttemptsCount >= MaxFailedLoginAttemptsPerIpAddress {
		return accessToken, errors.New("MaxFailedLoginAttemptsReached")
	}

	if !authQueryRepo.IsLoginValid(createDto) {
		recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
		NewCreateSecurityActivityRecord(activityRecordCmdRepo).
			CreateSessionToken(recordCode, createDto)

		return accessToken, errors.New("InvalidCredentials")
	}

	accountEntity, err := accountQueryRepo.ReadByUsername(createDto.Username)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	recordCode, _ := valueObject.NewActivityRecordCode("LoginSuccessful")
	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSessionToken(recordCode, createDto)

	expiresIn := valueObject.NewUnixTimeAfterNow(SessionTokenExpiresIn)

	return authCmdRepo.CreateSessionToken(
		accountEntity.Id, expiresIn, createDto.OperatorIpAddress,
	)
}
