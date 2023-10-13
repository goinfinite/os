package valueObject

import (
	"errors"
	"math/big"
)

type SslSerialNumber string

func NewSslSerialNumber(sslSerialNumber interface{}) (SslSerialNumber, error) {
	var sslSerialNumberOutput SslSerialNumber
	isValidType := true

	switch serialNumber := sslSerialNumber.(type) {
	case *big.Int:
		sslSerialNumberOutput = SslSerialNumber(serialNumber.String())
	case string:
		sslSerialNumberBigInt, err := stringToBigInt(serialNumber)
		if err != nil {
			isValidType = false
		}
		sslSerialNumberOutput = SslSerialNumber(sslSerialNumberBigInt.String())
	default:
		isValidType = false
	}

	if !isValidType || !sslSerialNumberOutput.isValid() {
		return "", errors.New("InvalidSslSerialNumber")
	}

	return sslSerialNumberOutput, nil
}

func NewSslSerialNumberPanic(input interface{}) SslSerialNumber {
	sslSerialNumber, err := NewSslSerialNumber(input)
	if err != nil {
		panic(err)
	}
	return sslSerialNumber
}

func stringToBigInt(sslSerialNumberStr string) (*big.Int, error) {
	sslSerialNumberBigInt := new(big.Int)
	sslSerialNumberBigInt, ok := sslSerialNumberBigInt.SetString(sslSerialNumberStr, 10)
	if !ok {
		return nil, errors.New("InvalidSslSerialNumber")
	}

	return sslSerialNumberBigInt, nil
}

func (sslSerialNumber SslSerialNumber) isValid() bool {
	sslSerialNumberBigInt, _ := stringToBigInt(string(sslSerialNumber))
	zeroBigInt := new(big.Int)
	result := sslSerialNumberBigInt.Cmp(zeroBigInt)
	if result <= 0 {
		return false
	}

	return true
}

func (sslSerialNumber SslSerialNumber) BigInt() big.Int {
	sslSerialNumberBigInt, _ := stringToBigInt(sslSerialNumber.String())
	return *sslSerialNumberBigInt
}

func (sslSerialNumber SslSerialNumber) String() string {
	return string(sslSerialNumber)
}
