package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	accountInfra "github.com/speedianet/os/src/infra/account"
	filesInfra "github.com/speedianet/os/src/infra/files"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetAccounts	 godoc
// @Summary      GetAccounts
// @Description  List accounts.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Account
// @Router       /v1/account/ [get]
func GetAccountsController(c echo.Context) error {
	accountsQueryRepo := accountInfra.AccQueryRepo{}
	accountsList, err := useCase.GetAccounts(accountsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, accountsList)
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
func CreateAccountController(c echo.Context) error {
	requiredParams := []string{"username", "password"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	createAccountDto := dto.NewCreateAccount(
		valueObject.NewUsernamePanic(requestBody["username"].(string)),
		valueObject.NewPasswordPanic(requestBody["password"].(string)),
	)

	accQueryRepo := accountInfra.AccQueryRepo{}
	accCmdRepo := accountInfra.AccCmdRepo{}
	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	err := useCase.CreateAccount(
		accQueryRepo,
		accCmdRepo,
		filesQueryRepo,
		filesCmdRepo,
		createAccountDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "AccountCreated")
}

// UpdateAccount godoc
// @Summary      UpdateAccount
// @Description  Update an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateDto 	  body dto.UpdateAccount  true   "Only id is required."
// @Success      200 {object} object{} "AccountUpdated message or NewKeyString"
// @Router       /v1/account/ [put]
func UpdateAccountController(c echo.Context) error {
	requestBody, _ := apiHelper.GetRequestBody(c)

	var accountIdPtr *valueObject.AccountId
	if requestBody["id"] != nil {
		accountId := valueObject.NewAccountIdPanic(requestBody["id"])
		accountIdPtr = &accountId
	}

	var usernamePtr *valueObject.Username
	if requestBody["username"] != nil {
		username := valueObject.NewUsernamePanic(requestBody["username"].(string))
		usernamePtr = &username
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

	updateDto := dto.NewUpdateAccount(
		accountIdPtr,
		usernamePtr,
		passPtr,
		shouldUpdateApiKeyPtr,
	)

	accQueryRepo := accountInfra.AccQueryRepo{}
	accCmdRepo := accountInfra.AccCmdRepo{}

	if updateDto.Password != nil {
		err := useCase.UpdateAccountPassword(
			accQueryRepo,
			accCmdRepo,
			updateDto,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusInternalServerError, err.Error(),
			)
		}
	}

	if updateDto.ShouldUpdateApiKey != nil && *updateDto.ShouldUpdateApiKey {
		newKey, err := useCase.UpdateAccountApiKey(
			accQueryRepo,
			accCmdRepo,
			updateDto,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
		}

		return apiHelper.ResponseWrapper(c, http.StatusOK, newKey)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountUpdated")
}

// DeleteAccount godoc
// @Summary      DeleteAccount
// @Description  Delete an account.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        accountId 	  path   string  true  "Account ID that will be deleted."
// @Success      200 {object} object{} "AccountDeleted"
// @Router       /v1/account/{accountId}/ [delete]
func DeleteAccountController(c echo.Context) error {
	accountId := valueObject.NewAccountIdPanic(c.Param("accountId"))

	accQueryRepo := accountInfra.AccQueryRepo{}
	accCmdRepo := accountInfra.AccCmdRepo{}

	err := useCase.DeleteAccount(
		accQueryRepo,
		accCmdRepo,
		accountId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "AccountDeleted")
}
