package uiPresenterHelper

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/presentation/liaison"
)

func ShouldEnableInitialSetup(accountLiaison *liaison.AccountLiaison) bool {
	accountsServiceResponse := accountLiaison.Read(map[string]interface{}{})
	if accountsServiceResponse.Status != liaison.Success {
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
