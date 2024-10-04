package apiController

import (
	"github.com/labstack/echo/v4"
	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController() *AccountController {
	return &AccountController{
		accountService: service.NewAccountService(),
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
// @Summary      CreateNewAccount
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
		return apiHelper.ServiceTokenResponseWrapper(
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
	requestBody := map[string]interface{}{
		"id": c.Param("accountId"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.accountService.Delete(requestBody),
	)
}
