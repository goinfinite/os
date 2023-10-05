package valueObject

import (
	"errors"
	"strconv"

	voHelper "github.com/speedianet/sam/src/domain/valueObject/helper"
)

type SslId uint64

func NewSslId(value interface{}) (SslId, error) {
	sslId, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidSslId")
	}

	return SslId(sslId), nil
}

func NewSslIdPanic(value interface{}) SslId {
	sslId, err := NewSslId(value)
	if err != nil {
		panic(err)
	}
	return sslId
}

func (id SslId) Get() uint64 {
	return uint64(id)
}

func (id SslId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
