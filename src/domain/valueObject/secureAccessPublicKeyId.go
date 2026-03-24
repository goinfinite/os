package valueObject

import (
	"errors"
	"strconv"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type SecureAccessPublicKeyId uint16

func NewSecureAccessPublicKeyId(value interface{}) (keyId SecureAccessPublicKeyId, err error) {
	uintValue, err := tkVoUtil.InterfaceToUint(value)
	if err != nil {
		return keyId, errors.New("SecureAccessPublicKeyIdMustBeUint")
	}

	return SecureAccessPublicKeyId(uintValue), nil
}

func (vo SecureAccessPublicKeyId) Uint16() uint16 {
	return uint16(vo)
}

func (vo SecureAccessPublicKeyId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
