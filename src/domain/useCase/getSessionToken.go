package useCase

import (
	"log"
	"time"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func GetSessionToken(
	authQueryRepo repository.AuthQueryRepo,
	authCmdRepo repository.AuthCmdRepo,
	accQueryRepo repository.AccQueryRepo,
	login dto.Login,
	ipAddress valueObject.IpAddress,
) entity.AccessToken {
	isLoginValid := authQueryRepo.IsLoginValid(login)

	if !isLoginValid {
		log.Printf(
			"Login failed for '%v' from '%v'.",
			login.Username.String(),
			ipAddress.String(),
		)
		panic("InvalidLoginCredentials")
	}

	accountDetails := accQueryRepo.GetAccountDetailsByUsername(login.Username)
	userId := accountDetails.UserId
	expiresIn := valueObject.UnixTime(
		time.Now().Add(3 * time.Hour).Unix(),
	)

	return authCmdRepo.GenerateSessionToken(userId, expiresIn, ipAddress)
}
