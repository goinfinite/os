package service

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	filesInfra "github.com/goinfinite/os/src/infra/files"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type AccountService struct {
}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (service *AccountService) Read() ServiceOutput {
	accountsQueryRepo := accountInfra.AccQueryRepo{}
	accountsList, err := useCase.GetAccounts(accountsQueryRepo)
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

	dto := dto.NewCreateAccount(username, password)

	accQueryRepo := accountInfra.AccQueryRepo{}
	accCmdRepo := accountInfra.AccCmdRepo{}
	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	err = useCase.CreateAccount(
		accQueryRepo, accCmdRepo, filesQueryRepo, filesCmdRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "AccountCreated")
}

func (service *AccountService) Update(input map[string]interface{}) ServiceOutput {
	var accountIdPtr *valueObject.AccountId
	if input["id"] != nil {
		accountId, err := valueObject.NewAccountId(input["id"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		accountIdPtr = &accountId
	}

	var usernamePtr *valueObject.Username
	if input["username"] != nil {
		username, err := valueObject.NewUsername(input["username"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		usernamePtr = &username
	}

	var passwordPtr *valueObject.Password
	if input["password"] != nil {
		password, err := valueObject.NewPassword(input["password"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		passwordPtr = &password
	}

	var shouldUpdateApiKeyPtr *bool
	if input["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey := input["shouldUpdateApiKey"].(bool)
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	dto := dto.NewUpdateAccount(
		accountIdPtr, usernamePtr, passwordPtr, shouldUpdateApiKeyPtr,
	)

	accQueryRepo := accountInfra.AccQueryRepo{}
	accCmdRepo := accountInfra.AccCmdRepo{}

	if dto.Password != nil {
		err := useCase.UpdateAccountPassword(accQueryRepo, accCmdRepo, dto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}
	}

	if dto.ShouldUpdateApiKey != nil && *dto.ShouldUpdateApiKey {
		newApiKey, err := useCase.UpdateAccountApiKey(accQueryRepo, accCmdRepo, dto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Success, newApiKey)
	}

	return NewServiceOutput(Success, "AccountUpdated")
}

func (service *AccountService) Delete(input map[string]interface{}) ServiceOutput {
	accountId, err := valueObject.NewAccountId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	accQueryRepo := accountInfra.AccQueryRepo{}
	accCmdRepo := accountInfra.AccCmdRepo{}

	err = useCase.DeleteAccount(accQueryRepo, accCmdRepo, accountId)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "AccountDeleted")
}
