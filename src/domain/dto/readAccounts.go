package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadAccountsRequest struct {
	Pagination                          tkDto.Pagination        `json:"pagination"`
	AccountId                           *tkValueObject.AccountId `json:"id,omitempty"`
	AccountUsername                     *valueObject.Username   `json:"username,omitempty"`
	ShouldIncludeSecureAccessPublicKeys *bool                   `json:"shouldIncludeSecureAccessPublicKeys,omitempty"`
}

type ReadAccountsResponse struct {
	Pagination tkDto.Pagination `json:"pagination"`
	Accounts   []entity.Account `json:"accounts"`
}
