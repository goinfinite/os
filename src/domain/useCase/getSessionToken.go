package useCase

import (
	"errors"
	"log"
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
	ipAddress valueObject.IpAddress,
) (entity.AccessToken, error) {
	isLoginValid := authQueryRepo.IsLoginValid(login)

	if !isLoginValid {
		log.Printf(
			"Login failed for '%v' from '%v'.",
			login.Username.String(),
			ipAddress.String(),
		)
		return entity.AccessToken{}, errors.New("InvalidCredentials")
	}

	accountDetails, err := accQueryRepo.GetByUsername(login.Username)
	if err != nil {
		return entity.AccessToken{}, errors.New("AccountNotFound")
	}

	accountId := accountDetails.Id
	expiresIn := valueObject.UnixTime(
		time.Now().Add(3 * time.Hour).Unix(),
	)

	return authCmdRepo.GenerateSessionToken(accountId, expiresIn, ipAddress)
}
