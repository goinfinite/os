package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	authInfra "github.com/goinfinite/os/src/infra/auth"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
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
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"username", "password", "operatorIpAddress"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	username, err := valueObject.NewUsername(untrustedInput["username"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	password, err := tkValueObject.NewWeakPassword(untrustedInput["password"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	dto := dto.NewCreateSessionToken(username, password, operatorIpAddress)

	authQueryRepo := authInfra.NewAuthQueryRepo(liaison.persistentDbSvc)
	authCmdRepo := authInfra.NewAuthCmdRepo()
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
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, accessToken,
	)
}
