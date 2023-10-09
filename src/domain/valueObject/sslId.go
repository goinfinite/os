package valueObject

import (
	"errors"
	"math/big"
)

type SslId struct {
	bigIntValue big.Int
	stringValue string
}

func NewSslId(value string) (SslId, error) {
	sslIdBigInt := new(big.Int)
	sslIdBigInt, ok := sslIdBigInt.SetString(value, 10)
	if !ok {
		return SslId{}, errors.New("InvalidSslId")
	}

	zeroBigInt := new(big.Int)
	result := sslIdBigInt.Cmp(zeroBigInt)
	if result <= 0 {
		return SslId{}, errors.New("InvalidSslId")
	}

	return SslId{
		bigIntValue: *sslIdBigInt,
		stringValue: sslIdBigInt.String(),
	}, nil
}

func NewSslIdPanic(value string) SslId {
	sslId, err := NewSslId(value)
	if err != nil {
		panic(err)
	}
	return sslId
}

func (id SslId) Get() big.Int {
	return id.bigIntValue
}

func (id SslId) String() string {
	return id.stringValue
}
