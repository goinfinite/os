package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	apiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

// GetAccounts	 godoc
// @Summary      GetAccounts
// @Description  List accounts.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.AccountDetails
// @Router       /account/ [get]
func GetAccountsController(c echo.Context) error {
	accountsQueryRepo := infra.AccQueryRepo{}
	accountsList, err := useCase.GetAccounts(accountsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accountsList)
}

// AddAccount    godoc
// @Summary      AddNewAccount
// @Description  Add a new account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addAccountDto 	  body    dto.AddAccount  true  "NewAccountDetails"
// @Success      201 {object} object{} "AccountCreated"
// @Router       /account/ [post]
func AddAccountController(c echo.Context) error {
	requiredParams := []string{"username", "password"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	addAccountDto := dto.NewAddAccount(
		valueObject.NewUsernamePanic(requestBody["username"].(string)),
		valueObject.NewPasswordPanic(requestBody["password"].(string)),
	)

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

	err := useCase.AddAccount(
		accQueryRepo,
		accCmdRepo,
		addAccountDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "AccountCreated")
}

// DeleteAccount godoc
// @Summary      DeleteAccount
// @Description  Delete an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "AccountId"
// @Success      200 {object} object{} "AccountDeleted"
// @Router       /account/{accountId}/ [delete]
func DeleteAccountController(c echo.Context) error {
	accountId := valueObject.NewAccountIdFromStringPanic(c.Param("accountId"))

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

	err := useCase.DeleteAccount(
		accQueryRepo,
		accCmdRepo,
		accountId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountDeleted")
}

// UpdateAccount godoc
// @Summary      UpdateAccount
// @Description  Update an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateAccountDto 	  body dto.UpdateAccount  true  "UpdateAccountDetails"
// @Success      200 {object} object{} "AccountUpdated message or NewKeyString"
// @Router       /account/ [put]
func UpdateAccountController(c echo.Context) error {
	requiredParams := []string{"accountId"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	var accountId valueObject.AccountId
	switch id := requestBody["accountId"].(type) {
	case string:
		accountId = valueObject.NewAccountIdFromStringPanic(id)
	case float64:
		accountId = valueObject.NewAccountIdFromFloatPanic(id)
	}

	var passPtr *valueObject.Password
	if requestBody["password"] != nil {
		password := valueObject.NewPasswordPanic(requestBody["password"].(string))
		passPtr = &password
	}

	var shouldUpdateApiKeyPtr *bool
	if requestBody["shouldUpdateApiKey"] != nil {
		shouldUpdateApiKey := requestBody["shouldUpdateApiKey"].(bool)
		shouldUpdateApiKeyPtr = &shouldUpdateApiKey
	}

	updateAccountDto := dto.NewUpdateAccount(
		accountId,
		passPtr,
		shouldUpdateApiKeyPtr,
	)

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

	if updateAccountDto.Password != nil {
		useCase.UpdateAccountPassword(
			accQueryRepo,
			accCmdRepo,
			updateAccountDto,
		)
	}

	if updateAccountDto.ShouldUpdateApiKey != nil && *updateAccountDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			accQueryRepo,
			accCmdRepo,
			updateAccountDto,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		return apiHelper.ResponseWrapper(c, http.StatusOK, newKey)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountUpdated")
}
