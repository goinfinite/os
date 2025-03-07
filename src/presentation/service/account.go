package service

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

var LocalOperatorAccountId, _ = valueObject.NewAccountId(0)
var LocalOperatorIpAddress = valueObject.NewLocalhostIpAddress()

type AccountService struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	accountQueryRepo      *accountInfra.AccountQueryRepo
	accountCmdRepo        *accountInfra.AccountCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewAccountService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountService {
	return &AccountService{
		persistentDbSvc:       persistentDbSvc,
		accountQueryRepo:      accountInfra.NewAccountQueryRepo(persistentDbSvc),
		accountCmdRepo:        accountInfra.NewAccountCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *AccountService) Read(input map[string]interface{}) ServiceOutput {
	if input["id"] != nil {
		input["accountId"] = input["id"]
	}

	var idPtr *valueObject.AccountId
	if input["id"] != nil {
		id, err := valueObject.NewAccountId(input["id"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		idPtr = &id
	}

	var usernamePtr *valueObject.Username
	if input["name"] != nil {
		username, err := valueObject.NewUsername(input["username"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		usernamePtr = &username
	}

	shouldIncludeSecureAccessPublicKeys := false
	if input["shouldIncludeSecureAccessPublicKeys"] != nil {
		var err error
		shouldIncludeSecureAccessPublicKeys, err = voHelper.InterfaceToBool(
			input["shouldIncludeSecureAccessPublicKeys"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			input["sortDirection"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readRequestDto := dto.ReadAccountsRequest{
		Pagination:                          paginationDto,
		AccountId:                           idPtr,
		AccountUsername:                     usernamePtr,
		ShouldIncludeSecureAccessPublicKeys: &shouldIncludeSecureAccessPublicKeys,
	}

	accountsList, err := useCase.ReadAccounts(service.accountQueryRepo, readRequestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, accountsList)
}

func (service *AccountService) Create(input map[string]interface{}) ServiceOutput {
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

	isSuperAdmin := false
	if input["isSuperAdmin"] != nil {
		isSuperAdmin, err = voHelper.InterfaceToBool(input["isSuperAdmin"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateAccount(
		username, password, isSuperAdmin, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateAccount(
		service.accountQueryRepo, service.accountCmdRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "AccountCreated")
}

func (service *AccountService) Update(input map[string]interface{}) ServiceOutput {
	if input["id"] != nil {
		input["accountId"] = input["id"]
	}

	requiredParams := []string{"accountId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var passwordPtr *valueObject.Password
	if input["password"] != nil && input["password"] != "" {
		password, err := valueObject.NewPassword(input["password"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		passwordPtr = &password
	}

	var isSuperAdminPtr *bool
	if input["isSuperAdmin"] != nil {
		isSuperAdmin, err := voHelper.InterfaceToBool(input["isSuperAdmin"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		isSuperAdminPtr = &isSuperAdmin
	}

	var shouldUpdateApiKeyPtr *bool
	if input["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey, err := voHelper.InterfaceToBool(input["shouldUpdateApiKey"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateAccount(
		accountId, passwordPtr, isSuperAdminPtr, shouldUpdateApiKeyPtr,
		operatorAccountId, operatorIpAddress,
	)

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			service.accountQueryRepo, service.accountCmdRepo,
			service.activityRecordCmdRepo, updateDto,
		)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}
		return NewServiceOutput(Success, newKey)
	}

	err = useCase.UpdateAccount(
		service.accountQueryRepo, service.accountCmdRepo,
		service.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "AccountUpdated")
}

func (service *AccountService) Delete(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"accountId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteAccount(accountId, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteAccount(
		service.accountQueryRepo, service.accountCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "AccountDeleted")
}

func (service *AccountService) CreateSecureAccessPublicKey(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"accountId", "content"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	keyContent, err := valueObject.NewSecureAccessPublicKeyContent(input["content"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	keyName, err := valueObject.NewSecureAccessPublicKeyName(input["name"])
	if err != nil {
		keyName, err = keyContent.ReadOnlyKeyName()
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	accountId, err := valueObject.NewAccountId(input["accountId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateSecureAccessPublicKey(
		accountId, keyContent, keyName, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateSecureAccessPublicKey(
		service.accountCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "SecureAccessPublicKeyCreated")
}

func (service *AccountService) DeleteSecureAccessPublicKey(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"secureAccessPublicKeyId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	keyId, err := valueObject.NewSecureAccessPublicKeyId(
		input["secureAccessPublicKeyId"],
	)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteSecureAccessPublicKey(
		keyId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteSecureAccessPublicKey(
		service.accountQueryRepo, service.accountCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "SecureAccessPublicKeyDeleted")
}
