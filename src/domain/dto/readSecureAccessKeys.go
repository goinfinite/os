package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadSecureAccessKeysRequest struct {
	Pagination          Pagination                       `json:"pagination"`
	AccountId           valueObject.AccountId            `json:"accountId,omitempty"`
	SecureAccessKeyId   *valueObject.SecureAccessKeyId   `json:"id,omitempty"`
	SecureAccessKeyName *valueObject.SecureAccessKeyName `json:"name,omitempty"`
}

type ReadSecureAccessKeysResponse struct {
	Pagination       Pagination               `json:"pagination"`
	SecureAccessKeys []entity.SecureAccessKey `json:"secureAccessKeys"`
}
