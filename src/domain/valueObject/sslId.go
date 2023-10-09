package valueObject

import (
	"errors"
	"math/big"
)

type SslId big.Int

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

	return SslId(*sslIdBigInt), nil
}

func NewSslIdPanic(value string) SslId {
	sslId, err := NewSslId(value)
	if err != nil {
		panic(err)
	}
	return sslId
}

func (id SslId) Get() big.Int {
	return big.Int(id)
}

func (id SslId) String() string {
	return big.NewInt(1).String()
}
