package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

var LocalOperatorAccountId = tkValueObject.AccountIdSystem
var LocalOperatorIpAddress = tkValueObject.IpAddressLocal

type AccountLiaison struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	accountQueryRepo      *accountInfra.AccountQueryRepo
	accountCmdRepo        *accountInfra.AccountCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewAccountLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountLiaison {
	return &AccountLiaison{
		persistentDbSvc:       persistentDbSvc,
		accountQueryRepo:      accountInfra.NewAccountQueryRepo(persistentDbSvc),
		accountCmdRepo:        accountInfra.NewAccountCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (liaison *AccountLiaison) Read(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	if untrustedInput["id"] != nil && untrustedInput["accountId"] == nil {
		untrustedInput["accountId"] = untrustedInput["id"]
	}

	var idPtr *tkValueObject.AccountId
	if untrustedInput["accountId"] != nil {
		id, err := tkValueObject.NewAccountId(untrustedInput["accountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err,
			)
		}
		idPtr = &id
	}

	var usernamePtr *valueObject.Username
	if untrustedInput["username"] != nil {
		username, err := valueObject.NewUsername(untrustedInput["username"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err,
			)
		}
		usernamePtr = &username
	}

	shouldIncludeSecureAccessPublicKeys := false
	if untrustedInput["shouldIncludeSecureAccessPublicKeys"] != nil {
		var err error
		shouldIncludeSecureAccessPublicKeys, err = tkVoUtil.InterfaceToBool(
			untrustedInput["shouldIncludeSecureAccessPublicKeys"],
		)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err,
			)
		}
	}

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.AccountsDefaultPagination, untrustedInput,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	readRequestDto := dto.ReadAccountsRequest{
		Pagination:                          requestPagination,
		AccountId:                           idPtr,
		AccountUsername:                     usernamePtr,
		ShouldIncludeSecureAccessPublicKeys: &shouldIncludeSecureAccessPublicKeys,
	}

	accountsList, err := useCase.ReadAccounts(liaison.accountQueryRepo, readRequestDto)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, accountsList,
	)
}

func (liaison *AccountLiaison) Create(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"username", "password"}
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

	password, err := tkValueObject.NewPassword(untrustedInput["password"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	isSuperAdmin := false
	if untrustedInput["isSuperAdmin"] != nil {
		isSuperAdmin, err = tkVoUtil.InterfaceToBool(untrustedInput["isSuperAdmin"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	createDto := dto.NewCreateAccount(
		username, password, isSuperAdmin, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateAccount(
		liaison.accountQueryRepo, liaison.accountCmdRepo,
		liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusCreated, "AccountCreated",
	)
}

func (liaison *AccountLiaison) Update(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	if untrustedInput["id"] != nil {
		untrustedInput["accountId"] = untrustedInput["id"]
	}

	if untrustedInput["username"] != nil {
		untrustedInput["accountUsername"] = untrustedInput["username"]
	}

	var err error
	var accountIdPtr *tkValueObject.AccountId
	if untrustedInput["accountId"] != nil {
		accountId, err := tkValueObject.NewAccountId(untrustedInput["accountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		accountIdPtr = &accountId
	}

	var accountUsernamePtr *valueObject.Username
	if untrustedInput["accountUsername"] != nil {
		accountUsername, err := valueObject.NewUsername(untrustedInput["accountUsername"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		accountUsernamePtr = &accountUsername
	}

	var passwordPtr *tkValueObject.Password
	if untrustedInput["password"] != nil && untrustedInput["password"] != "" {
		password, err := tkValueObject.NewPassword(untrustedInput["password"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		passwordPtr = &password
	}

	var isSuperAdminPtr *bool
	if untrustedInput["isSuperAdmin"] != nil {
		isSuperAdmin, err := tkVoUtil.InterfaceToBool(untrustedInput["isSuperAdmin"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		isSuperAdminPtr = &isSuperAdmin
	}

	var shouldUpdateApiKeyPtr *bool
	if untrustedInput["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey, err := tkVoUtil.InterfaceToBool(untrustedInput["shouldUpdateApiKey"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	updateDto := dto.NewUpdateAccount(
		accountIdPtr, accountUsernamePtr, passwordPtr, isSuperAdminPtr, shouldUpdateApiKeyPtr,
		operatorAccountId, operatorIpAddress,
	)

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			liaison.accountQueryRepo, liaison.accountCmdRepo,
			liaison.activityRecordCmdRepo, updateDto,
		)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
			)
		}
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusSuccess, newKey,
		)
	}

	err = useCase.UpdateAccount(
		liaison.accountQueryRepo, liaison.accountCmdRepo,
		liaison.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, "AccountUpdated",
	)
}

func (liaison *AccountLiaison) Delete(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	if untrustedInput["id"] != nil && untrustedInput["accountId"] == nil {
		untrustedInput["accountId"] = untrustedInput["id"]
	}

	requiredParams := []string{"accountId"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	accountId, err := tkValueObject.NewAccountId(untrustedInput["accountId"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	deleteDto := dto.NewDeleteAccount(accountId, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteAccount(
		liaison.accountQueryRepo, liaison.accountCmdRepo,
		liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, "AccountDeleted",
	)
}

func (liaison *AccountLiaison) CreateSecureAccessPublicKey(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"accountId", "content"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	keyContent, err := valueObject.NewSecureAccessPublicKeyContent(untrustedInput["content"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	keyName, err := valueObject.NewSecureAccessPublicKeyName(untrustedInput["name"])
	if err != nil {
		keyName, err = keyContent.ReadOnlyKeyName()
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	accountId, err := tkValueObject.NewAccountId(untrustedInput["accountId"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	createDto := dto.NewCreateSecureAccessPublicKey(
		accountId, keyContent, keyName, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateSecureAccessPublicKey(
		liaison.accountCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusCreated, "SecureAccessPublicKeyCreated",
	)
}

func (liaison *AccountLiaison) DeleteSecureAccessPublicKey(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"secureAccessPublicKeyId"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	keyId, err := valueObject.NewSecureAccessPublicKeyId(
		untrustedInput["secureAccessPublicKeyId"],
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusUserError, err.Error(),
		)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(
				tkPresentation.LiaisonResponseStatusUserError, err.Error(),
			)
		}
	}

	deleteDto := dto.NewDeleteSecureAccessPublicKey(
		keyId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteSecureAccessPublicKey(
		liaison.accountQueryRepo, liaison.accountCmdRepo,
		liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(
			tkPresentation.LiaisonResponseStatusInfraError, err.Error(),
		)
	}

	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusSuccess, "SecureAccessPublicKeyDeleted",
	)
}
