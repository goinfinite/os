package useCase

import (
	"log"

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
	return authCmdRepo.GenerateSessionToken(userId, 28800, ipAddress)
}
