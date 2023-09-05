package valueObject

import (
	"errors"
	"reflect"
	"strconv"
)

type CronId uint64

func NewCronId(value interface{}) (CronId, error) {
	var cronId uint64
	var err error
	switch v := value.(type) {
	case string:
		cronId, err = strconv.ParseUint(v, 10, 64)
	case int, int8, int16, int32, int64:
		intValue := reflect.ValueOf(v).Int()
		if intValue < 0 {
			err = errors.New("InvalidCronId")
		}
		cronId = uint64(intValue)
	case uint, uint8, uint16, uint32, uint64:
		cronId = uint64(reflect.ValueOf(v).Uint())
	case float32, float64:
		floatValue := reflect.ValueOf(v).Float()
		if floatValue < 0 {
			err = errors.New("InvalidCronId")
		}
		cronId = uint64(floatValue)
	default:
		err = errors.New("InvalidCronId")
	}

	if err != nil {
		return 0, errors.New("InvalidCronId")
	}

	return CronId(cronId), nil
}

func NewCronIdPanic(value interface{}) CronId {
	cronId, err := NewCronId(value)
	if err != nil {
		panic(err)
	}
	return cronId
}

func (id CronId) Get() uint64 {
	return uint64(id)
}

func (id CronId) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
