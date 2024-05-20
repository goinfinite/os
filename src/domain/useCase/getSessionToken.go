package useCase

import (
	"errors"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func GetSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	login dto.Login,
) (accessToken entity.AccessToken, err error) {
	if !authQueryRepo.IsLoginValid(login) {
		return accessToken, errors.New("InvalidCredentials")
	}

	accountDetails, err := accQueryRepo.GetByUsername(login.Username)
	if err != nil {
		return accessToken, errors.New("AccountNotFound")
	}

	accountId := accountDetails.Id
	expiresIn := valueObject.UnixTime(time.Now().Add(3 * time.Hour).Unix())

	return authCmdRepo.GenerateSessionToken(accountId, expiresIn, login.IpAddress)
}
