package apiController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/labstack/echo/v4"
)

type AccountController struct {
	accountLiaison *liaison.AccountLiaison
}

func NewAccountController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountController {
	return &AccountController{
		accountLiaison: liaison.NewAccountLiaison(persistentDbSvc, trailDbSvc),
	}
}

// ReadAccounts	 godoc
// @Summary      ReadAccounts
// @Description  List accounts.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id query  string  false  "Id"
// @Param        username query  string  false  "Username"
// @Param        shouldIncludeSecureAccessPublicKeys query  bool  false  "ShouldIncludeSecureAccessPublicKeys (only works if OpenSSH service is installed)"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadAccountsResponse
// @Router       /v1/account/ [get]
func (controller *AccountController) Read(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.accountLiaison.Read(requestInputData),
	)
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.accountLiaison.Create(requestInputData),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.accountLiaison.Update(requestInputData),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.accountLiaison.Delete(requestInputData),
	)
}

// CreateSecureAccessPublicKey    godoc
// @Summary      CreateSecureAccessPublicKey
// @Description  Create a new secure access public key.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createSecureAccessPublicKey 	  body    dto.CreateSecureAccessPublicKey  true  "'name' is optional. Will only become required if there is no name in 'content'. If the 'name' is provided, it will overwrite the name in the 'content'."
// @Success      201 {object} object{} "SecureAccessPublicKeyCreated"
// @Router       /v1/account/secure-access-public-key/ [post]
func (controller *AccountController) CreateSecureAccessPublicKey(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.accountLiaison.CreateSecureAccessPublicKey(requestInputData),
	)
}

// DeleteSecureAccessPublicKey godoc
// @Summary      DeleteSecureAccessPublicKey
// @Description  Delete a secure access public key.
// @Tags         account
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        secureAccessPublicKeyId 	  path   string  true  "SecureAccessPublicKeyId to delete."
// @Success      200 {object} object{} "SecureAccessPublicKeyDeleted"
// @Router       /v1/account/secure-access-public-key/{secureAccessPublicKeyId}/ [delete]
func (controller *AccountController) DeleteSecureAccessPublicKey(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.accountLiaison.DeleteSecureAccessPublicKey(requestInputData),
	)
}
