package valueObject

import (
	"errors"
	"math/big"
)

type SslSerialNumber struct {
	bigIntValue big.Int
	stringValue string
}

func NewSslSerialNumber(value string) (SslSerialNumber, error) {
	sslSerialNumberBigInt := new(big.Int)
	sslSerialNumberBigInt, ok := sslSerialNumberBigInt.SetString(value, 10)
	if !ok {
		return SslSerialNumber{}, errors.New("InvalidSslSerialNumber")
	}

	zeroBigInt := new(big.Int)
	result := sslSerialNumberBigInt.Cmp(zeroBigInt)
	if result <= 0 {
		return SslSerialNumber{}, errors.New("InvalidSslSerialNumber")
	}

	return SslSerialNumber{
		bigIntValue: *sslSerialNumberBigInt,
		stringValue: sslSerialNumberBigInt.String(),
	}, nil
}

func NewSslSerialNumberPanic(value string) SslSerialNumber {
	sslSerialNumber, err := NewSslSerialNumber(value)
	if err != nil {
		panic(err)
	}
	return sslSerialNumber
}

func (id SslSerialNumber) Get() big.Int {
	return id.bigIntValue
}

func (id SslSerialNumber) String() string {
	return id.stringValue
}
