package useCase

import (
	"errors"
	"time"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

const MaxFailedLoginAttemptsPerIpAddress uint = 3
const FailedLoginAttemptsInterval time.Duration = 15 * time.Minute
const SessionTokenExpiresIn time.Duration = 3 * time.Hour

func getFailedLoginAttemptsCount(
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	loginDto dto.Login,
) uint {
	secLevel, _ := valueObject.NewActivityRecordLevel("SEC")
	recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
	failedAttemptsIntervalStartsAt := valueObject.NewUnixTimeBeforeNow(
		FailedLoginAttemptsInterval,
	)
	readActivityRecordsDto := dto.NewReadActivityRecords(
		&secLevel, &recordCode, nil, &loginDto.IpAddress, nil, nil, nil,
		nil, &failedAttemptsIntervalStartsAt,
	)

	failedLoginAttempts := ReadActivityRecords(
		activityRecordQueryRepo, readActivityRecordsDto,
	)

	return uint(len(failedLoginAttempts))
}

func GetSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	activityRecordQueryRepo repository.ActivityRecordQueryRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	loginDto dto.Login,
) (accessToken entity.AccessToken, err error) {
	failedAttemptsCount := getFailedLoginAttemptsCount(
		activityRecordQueryRepo, loginDto,
	)
	if failedAttemptsCount >= MaxFailedLoginAttemptsPerIpAddress {
		return accessToken, errors.New("MaxFailedLoginAttemptsReached")
	}

	if !authQueryRepo.IsLoginValid(loginDto) {
		recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
		CreateSecurityActivityRecord(
			activityRecordCmdRepo, &recordCode, &loginDto.IpAddress,
			nil, nil, &loginDto.Username,
		)

		return accessToken, errors.New("InvalidCredentials")
	}

	accountDetails, err := accQueryRepo.GetByUsername(loginDto.Username)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	accountId := accountDetails.Id
	recordCode, _ := valueObject.NewActivityRecordCode("LoginSuccessful")
	CreateSecurityActivityRecord(
		activityRecordCmdRepo, &recordCode, &loginDto.IpAddress, nil,
		&accountId, &loginDto.Username,
	)

	expiresIn := valueObject.NewUnixTimeAfterNow(SessionTokenExpiresIn)

	return authCmdRepo.GenerateSessionToken(accountId, expiresIn, loginDto.IpAddress)
}
