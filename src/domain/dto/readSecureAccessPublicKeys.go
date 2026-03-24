package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadSecureAccessPublicKeysRequest struct {
	Pagination                tkDto.Pagination                       `json:"pagination"`
	AccountId                 tkValueObject.AccountId                `json:"accountId,omitempty"`
	SecureAccessPublicKeyId   *valueObject.SecureAccessPublicKeyId   `json:"id,omitempty"`
	SecureAccessPublicKeyName *valueObject.SecureAccessPublicKeyName `json:"name,omitempty"`
}

type ReadSecureAccessPublicKeysResponse struct {
	Pagination             tkDto.Pagination               `json:"pagination"`
	SecureAccessPublicKeys []entity.SecureAccessPublicKey `json:"SecureAccessPublicKeys"`
}
