package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	accountInfra "github.com/speedianet/os/src/infra/account"
	activityRecordInfra "github.com/speedianet/os/src/infra/activityRecord"
	authInfra "github.com/speedianet/os/src/infra/auth"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
)

type AuthService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewAuthService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthService {
	return &AuthService{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (service *AuthService) GenerateJwtWithCredentials(
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

	dto := dto.NewCreateSessionToken(username, password, ipAddress)

	authQueryRepo := authInfra.AuthQueryRepo{}
	authCmdRepo := authInfra.AuthCmdRepo{}
	accountQueryRepo := accountInfra.NewAccountQueryRepo(service.persistentDbSvc)
	activityRecordQueryRepo := activityRecordInfra.NewActivityRecordQueryRepo(
		service.trailDbSvc,
	)
	activityRecordCmdRepo := activityRecordInfra.NewActivityRecordCmdRepo(
		service.trailDbSvc,
	)

	accessToken, err := useCase.CreateSessionToken(
		authQueryRepo, authCmdRepo, accountQueryRepo, activityRecordQueryRepo,
		activityRecordCmdRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, accessToken)
}
