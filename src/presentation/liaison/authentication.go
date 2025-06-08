package liaison

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	authInfra "github.com/goinfinite/os/src/infra/auth"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
)

type AuthenticationLiaison struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewAuthenticationLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthenticationLiaison {
	return &AuthenticationLiaison{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (liaison *AuthenticationLiaison) Login(
	untrustedInput map[string]any,
) LiaisonOutput {
	requiredParams := []string{"username", "password", "operatorIpAddress"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	username, err := valueObject.NewUsername(untrustedInput["username"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	password, err := valueObject.NewPassword(untrustedInput["password"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dto := dto.NewCreateSessionToken(username, password, operatorIpAddress)

	authQueryRepo := authInfra.NewAuthQueryRepo(liaison.persistentDbSvc)
	authCmdRepo := authInfra.AuthCmdRepo{}
	accountQueryRepo := accountInfra.NewAccountQueryRepo(liaison.persistentDbSvc)
	activityRecordQueryRepo := activityRecordInfra.NewActivityRecordQueryRepo(
		liaison.trailDbSvc,
	)
	activityRecordCmdRepo := activityRecordInfra.NewActivityRecordCmdRepo(
		liaison.trailDbSvc,
	)

	accessToken, err := useCase.CreateSessionToken(
		authQueryRepo, authCmdRepo, accountQueryRepo, activityRecordQueryRepo,
		activityRecordCmdRepo, dto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, accessToken)
}
