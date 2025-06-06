package liaison

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
)

var LocalOperatorAccountId = valueObject.AccountIdSystem
var LocalOperatorIpAddress = valueObject.IpAddressSystem

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

func (liaison *AccountLiaison) Read(untrustedInput map[string]any) LiaisonOutput {
	if untrustedInput["id"] != nil {
		untrustedInput["accountId"] = untrustedInput["id"]
	}

	var idPtr *valueObject.AccountId
	if untrustedInput["id"] != nil {
		id, err := valueObject.NewAccountId(untrustedInput["id"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		idPtr = &id
	}

	var usernamePtr *valueObject.Username
	if untrustedInput["name"] != nil {
		username, err := valueObject.NewUsername(untrustedInput["username"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		usernamePtr = &username
	}

	shouldIncludeSecureAccessPublicKeys := false
	if untrustedInput["shouldIncludeSecureAccessPublicKeys"] != nil {
		var err error
		shouldIncludeSecureAccessPublicKeys, err = voHelper.InterfaceToBool(
			untrustedInput["shouldIncludeSecureAccessPublicKeys"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
	}

	paginationDto := useCase.AccountsDefaultPagination
	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			untrustedInput["sortDirection"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readRequestDto := dto.ReadAccountsRequest{
		Pagination:                          paginationDto,
		AccountId:                           idPtr,
		AccountUsername:                     usernamePtr,
		ShouldIncludeSecureAccessPublicKeys: &shouldIncludeSecureAccessPublicKeys,
	}

	accountsList, err := useCase.ReadAccounts(liaison.accountQueryRepo, readRequestDto)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, accountsList)
}

func (liaison *AccountLiaison) Create(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"username", "password"}
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

	isSuperAdmin := false
	if untrustedInput["isSuperAdmin"] != nil {
		isSuperAdmin, err = voHelper.InterfaceToBool(untrustedInput["isSuperAdmin"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
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
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "AccountCreated")
}

func (liaison *AccountLiaison) Update(untrustedInput map[string]any) LiaisonOutput {
	if untrustedInput["id"] != nil {
		untrustedInput["accountId"] = untrustedInput["id"]
	}

	requiredParams := []string{"accountId"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(untrustedInput["accountId"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var passwordPtr *valueObject.Password
	if untrustedInput["password"] != nil && untrustedInput["password"] != "" {
		password, err := valueObject.NewPassword(untrustedInput["password"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		passwordPtr = &password
	}

	var isSuperAdminPtr *bool
	if untrustedInput["isSuperAdmin"] != nil {
		isSuperAdmin, err := voHelper.InterfaceToBool(untrustedInput["isSuperAdmin"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		isSuperAdminPtr = &isSuperAdmin
	}

	var shouldUpdateApiKeyPtr *bool
	if untrustedInput["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey, err := voHelper.InterfaceToBool(untrustedInput["shouldUpdateApiKey"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateAccount(
		accountId, passwordPtr, isSuperAdminPtr, shouldUpdateApiKeyPtr,
		operatorAccountId, operatorIpAddress,
	)

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			liaison.accountQueryRepo, liaison.accountCmdRepo,
			liaison.activityRecordCmdRepo, updateDto,
		)
		if err != nil {
			return NewLiaisonOutput(InfraError, err.Error())
		}
		return NewLiaisonOutput(Success, newKey)
	}

	err = useCase.UpdateAccount(
		liaison.accountQueryRepo, liaison.accountCmdRepo,
		liaison.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "AccountUpdated")
}

func (liaison *AccountLiaison) Delete(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"accountId"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(untrustedInput["accountId"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteAccount(accountId, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteAccount(
		liaison.accountQueryRepo, liaison.accountCmdRepo,
		liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "AccountDeleted")
}

func (liaison *AccountLiaison) CreateSecureAccessPublicKey(
	untrustedInput map[string]any,
) LiaisonOutput {
	requiredParams := []string{"accountId", "content"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	keyContent, err := valueObject.NewSecureAccessPublicKeyContent(untrustedInput["content"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	keyName, err := valueObject.NewSecureAccessPublicKeyName(untrustedInput["name"])
	if err != nil {
		keyName, err = keyContent.ReadOnlyKeyName()
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	accountId, err := valueObject.NewAccountId(untrustedInput["accountId"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateSecureAccessPublicKey(
		accountId, keyContent, keyName, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateSecureAccessPublicKey(
		liaison.accountCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "SecureAccessPublicKeyCreated")
}

func (liaison *AccountLiaison) DeleteSecureAccessPublicKey(
	untrustedInput map[string]any,
) LiaisonOutput {
	requiredParams := []string{"secureAccessPublicKeyId"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	keyId, err := valueObject.NewSecureAccessPublicKeyId(
		untrustedInput["secureAccessPublicKeyId"],
	)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
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
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "SecureAccessPublicKeyDeleted")
}
