package service

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	authInfra "github.com/goinfinite/os/src/infra/auth"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type AuthenticationService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewAuthenticationService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthenticationService {
	return &AuthenticationService{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (service *AuthenticationService) Login(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"username", "password"}
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

	operatorIpAddress, err := valueObject.NewIpAddress(input["operatorIpAddress"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dto := dto.NewCreateSessionToken(username, password, operatorIpAddress)

	authQueryRepo := authInfra.NewAuthQueryRepo(service.persistentDbSvc)
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
