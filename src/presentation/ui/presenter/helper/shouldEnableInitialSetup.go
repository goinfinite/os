package uiPresenterHelper

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/presentation/service"
)

func ShouldEnableInitialSetup(accountService *service.AccountService) bool {
	accountsServiceResponse := accountService.Read(map[string]interface{}{})
	if accountsServiceResponse.Status != service.Success {
		return false
	}

	accountsReadResponse, assertOk := accountsServiceResponse.Body.(dto.ReadAccountsResponse)
	if !assertOk {
		return false
	}

	if len(accountsReadResponse.Accounts) > 0 {
		return false
	}

	return true
}
