package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadSecureAccessPublicKeysRequest struct {
	Pagination                Pagination                             `json:"pagination"`
	AccountId                 valueObject.AccountId                  `json:"accountId,omitempty"`
	SecureAccessPublicKeyId   *valueObject.SecureAccessPublicKeyId   `json:"id,omitempty"`
	SecureAccessPublicKeyName *valueObject.SecureAccessPublicKeyName `json:"name,omitempty"`
}

type ReadSecureAccessPublicKeysResponse struct {
	Pagination             Pagination                     `json:"pagination"`
	SecureAccessPublicKeys []entity.SecureAccessPublicKey `json:"SecureAccessPublicKeys"`
}
