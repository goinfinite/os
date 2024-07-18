package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	accountInfra "github.com/speedianet/os/src/infra/account"
	authInfra "github.com/speedianet/os/src/infra/auth"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
)

type AuthService struct {
}

func NewAuthService() AuthService {
	return AuthService{}
}

func (service AuthService) GenerateJwtWithCredentials(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"username", "password", "ipAddress"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	username, err := valueObject.NewUsername(input["username"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	password, err := valueObject.NewPassword(input["password"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	ipAddress, err := valueObject.NewIpAddress(input["ipAddress"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	loginDto := dto.NewLogin(username, password, ipAddress)

	authQueryRepo := authInfra.AuthQueryRepo{}
	authCmdRepo := authInfra.AuthCmdRepo{}
	accQueryRepo := accountInfra.AccQueryRepo{}

	accessToken, err := useCase.GetSessionToken(
		authQueryRepo,
		authCmdRepo,
		accQueryRepo,
		loginDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, accessToken)
}
