package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type SecureAccessKeyId uint16

func NewSecureAccessKeyId(value interface{}) (keyId SecureAccessKeyId, err error) {
	uintValue, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return keyId, errors.New("SecureAccessKeyIdMustBeUint")
	}

	return SecureAccessKeyId(uintValue), nil
}

func (vo SecureAccessKeyId) Uint16() uint16 {
	return uint16(vo)
}

func (vo SecureAccessKeyId) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
