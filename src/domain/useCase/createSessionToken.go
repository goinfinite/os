package useCase

import (
	"errors"
	"time"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

const MaxFailedLoginAttemptsPerIpAddress uint = 3
const FailedLoginAttemptsInterval time.Duration = 15 * time.Minute
const SessionTokenExpiresIn time.Duration = 3 * time.Hour

func getFailedLoginAttemptsCount(
	activityRecordQueryRepo tkRepository.ActivityRecordQueryRepo,
	createDto dto.CreateSessionToken,
) uint {
	recordLevel := tkValueObject.ActivityRecordLevelSecurity
	recordCode, _ := tkValueObject.NewActivityRecordCode("LoginFailed")
	failedAttemptsIntervalStartsAt := tkValueObject.NewUnixTimeBeforeNow(
		FailedLoginAttemptsInterval,
	)
	readDto := tkDto.ReadActivityRecordsRequest{
		Pagination:        tkDto.PaginationUnpaginated,
		RecordLevel:       &recordLevel,
		RecordCode:        &recordCode,
		OperatorIpAddress: &createDto.OperatorIpAddress,
		CreatedAfterAt:    &failedAttemptsIntervalStartsAt,
	}

	responseDto, _ := ReadActivityRecords(activityRecordQueryRepo, readDto)
	return uint(len(responseDto.ActivityRecords))
}

func CreateSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accountQueryRepo repository.AccountQueryRepo,
	activityRecordQueryRepo tkRepository.ActivityRecordQueryRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	createDto dto.CreateSessionToken,
) (accessToken entity.AccessToken, err error) {
	failedAttemptsCount := getFailedLoginAttemptsCount(activityRecordQueryRepo, createDto)
	if failedAttemptsCount >= MaxFailedLoginAttemptsPerIpAddress {
		return accessToken, errors.New("MaxFailedLoginAttemptsReached")
	}

	if !authQueryRepo.IsLoginValid(createDto) {
		recordCode, _ := tkValueObject.NewActivityRecordCode("LoginFailed")
		NewCreateSecurityActivityRecord(activityRecordCmdRepo).
			CreateSessionToken(recordCode, createDto)

		return accessToken, errors.New("InvalidCredentials")
	}

	readRequestDto := dto.ReadAccountsRequest{
		AccountUsername: &createDto.Username,
	}
	accountEntity, err := accountQueryRepo.ReadFirst(readRequestDto)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	recordCode, _ := tkValueObject.NewActivityRecordCode("LoginSuccessful")
	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateSessionToken(recordCode, createDto)

	expiresIn := tkValueObject.NewUnixTimeAfterNow(SessionTokenExpiresIn)

	return authCmdRepo.CreateSessionToken(
		accountEntity.Id, expiresIn, createDto.OperatorIpAddress,
	)
}
