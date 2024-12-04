package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadAccountsRequest struct {
	Pagination                          Pagination             `json:"pagination"`
	AccountId                           *valueObject.AccountId `json:"id,omitempty"`
	AccountUsername                     *valueObject.Username  `json:"username,omitempty"`
	ShouldIncludeSecureAccessPublicKeys *bool                  `json:"shouldIncludeSecureAccessPublicKeys,omitempty"`
}

type ReadAccountsResponse struct {
	Pagination Pagination       `json:"pagination"`
	Accounts   []entity.Account `json:"accounts"`
}
