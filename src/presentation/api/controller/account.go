package apiController

import (
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountController {
	return &AccountController{
		accountService: service.NewAccountService(persistentDbSvc, trailDbSvc),
	}
}

// ReadAccounts	 godoc
// @Summary      ReadAccounts
// @Description  List accounts.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Account
// @Router       /v1/account/ [get]
func (controller *AccountController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.accountService.Read())
}

// CreateAccount    godoc
// @Summary      CreateAccount
// @Description  Create a new account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createAccountDto 	  body    dto.CreateAccount  true  "All props are required."
// @Success      201 {object} object{} "AccountCreated"
// @Router       /v1/account/ [post]
func (controller *AccountController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Create(requestBody),
	)
}

// UpdateAccount godoc
// @Summary      UpdateAccount
// @Description  Update an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateDto 	  body dto.UpdateAccount  true   "Only id or username is required."
// @Success      200 {object} object{} "'AccountUpdated' message or new API key in string format"
// @Router       /v1/account/ [put]
func (controller *AccountController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	shouldUpdateApiKey, err := voHelper.InterfaceToBool(
		requestBody["shouldUpdateApiKey"],
	)
	if err == nil && shouldUpdateApiKey {
		return apiHelper.ServiceResponseWithIgnoreToastHeaderWrapper(
			c, controller.accountService.Update(requestBody),
		)
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Update(requestBody),
	)
}

// DeleteAccount godoc
// @Summary      DeleteAccount
// @Description  Delete an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId to delete."
// @Success      200 {object} object{} "AccountDeleted"
// @Router       /v1/account/{accountId}/ [delete]
func (controller *AccountController) Delete(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Delete(requestBody),
	)
}

// ReadSecureAccessKeys	 godoc
// @Summary      ReadSecureAccessKeys
// @Description  List accounts secure access keys.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.SecureAccessKey
// @Router       /v1/account/secure-access-key/ [get]
func (controller *AccountController) ReadSecureAccessKey(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.ReadSecureAccessKey(requestBody),
	)
}

// CreateSecureAccessKey    godoc
// @Summary      CreateSecureAccessKey
// @Description  Create a new secure access key.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createSecureAccessKey 	  body    dto.CreateSecureAccessKey  true  "All props are required."
// @Success      201 {object} object{} "SecureAccessKeyCreated"
// @Router       /v1/account/secure-access-key/ [post]
func (controller *AccountController) CreateSecureAccessKey(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.CreateSecureAccessKey(requestBody),
	)
}
