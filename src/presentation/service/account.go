package service

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	secureAccessKeyInfra "github.com/goinfinite/os/src/infra/account/secureAccessKey"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
)

var LocalOperatorAccountId, _ = valueObject.NewAccountId(0)
var LocalOperatorIpAddress = valueObject.NewLocalhostIpAddress()

type AccountService struct {
	persistentDbSvc          *internalDbInfra.PersistentDatabaseService
	accountQueryRepo         *accountInfra.AccountQueryRepo
	accountCmdRepo           *accountInfra.AccountCmdRepo
	secureAccessKeyQueryRepo *secureAccessKeyInfra.SecureAccessKeyQueryRepo
	secureAccessKeyCmdRepo   *secureAccessKeyInfra.SecureAccessKeyCmdRepo
	activityRecordCmdRepo    *activityRecordInfra.ActivityRecordCmdRepo
	availabilityInspector    *sharedHelper.ServiceAvailabilityInspector
}

func NewAccountService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountService {
	return &AccountService{
		persistentDbSvc:  persistentDbSvc,
		accountQueryRepo: accountInfra.NewAccountQueryRepo(persistentDbSvc),
		accountCmdRepo:   accountInfra.NewAccountCmdRepo(persistentDbSvc),
		secureAccessKeyQueryRepo: secureAccessKeyInfra.NewSecureAccessKeyQueryRepo(
			persistentDbSvc,
		),
		secureAccessKeyCmdRepo: secureAccessKeyInfra.NewSecureAccessKeyCmdRepo(
			persistentDbSvc,
		),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
		availabilityInspector: sharedHelper.NewServiceAvailabilityInspector(
			persistentDbSvc,
		),
	}
}

func (service *AccountService) Read() ServiceOutput {
	accountsList, err := useCase.ReadAccounts(service.accountQueryRepo)
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
		username, password, operatorAccountId, operatorIpAddress,
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
	if input["password"] != nil {
		password, err := valueObject.NewPassword(input["password"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		passwordPtr = &password
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
		accountId, passwordPtr, shouldUpdateApiKeyPtr, operatorAccountId,
		operatorIpAddress,
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

func (service *AccountService) ReadSecureAccessKey(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("openssh")
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	if input["id"] != nil {
		input["accountId"] = input["id"]
	}

	if input["accountId"] == nil {
		input["accountId"] = input["operatorAccountId"]
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

	secureAccessKeys, err := useCase.ReadSecureAccessKeys(
		service.accountQueryRepo, service.secureAccessKeyQueryRepo, accountId,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, secureAccessKeys)
}

func (service *AccountService) CreateSecureAccessKey(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("openssh")
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	requiredParams := []string{"accountId", "content"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	keyContent, err := valueObject.NewSecureAccessKeyContent(input["content"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	keyName, err := valueObject.NewSecureAccessKeyName(input["name"])
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

	createDto := dto.NewCreateSecureAccessKey(
		keyName, keyContent, accountId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateSecureAccessKey(
		service.accountQueryRepo, service.secureAccessKeyCmdRepo,
		service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "SecureAccessKeyCreated")
}

func (service *AccountService) DeleteSecureAccessKey(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("openssh")
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	requiredParams := []string{"accountId", "secureAccessKeyId"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	keyId, err := valueObject.NewSecureAccessKeyId(input["secureAccessKeyId"])
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

	deleteDto := dto.NewDeleteSecureAccessKey(
		keyId, accountId, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteSecureAccessKey(
		service.accountQueryRepo, service.secureAccessKeyCmdRepo,
		service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "SecureAccessKeyDeleted")
}
