package valueObject

import (
	"errors"
	"reflect"
	"strconv"
)

type SslId uint64

func NewSslId(value interface{}) (SslId, error) {
	var sslId uint64
	var err error
	switch v := value.(type) {
	case string:
		sslId, err = strconv.ParseUint(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 {
			err = errors.New("InvalidSslId")
		}
		sslId = uint64(intValue)
	case uint, uint8, uint16, uint32, uint64:
		sslId = uint64(reflect.ValueOf(v).Uint())
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 {
			err = errors.New("InvalidSslId")
		}
		sslId = uint64(floatValue)
	default:
		err = errors.New("InvalidSslId")
	}

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
