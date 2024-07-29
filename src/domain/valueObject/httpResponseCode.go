package valueObject

import (
	"errors"
	"regexp"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const responseCodeExpression = "^([1-5][0-9][0-9])$"

type HttpResponseCode uint64

func NewHttpResponseCode(value interface{}) (
	httpResponseCode HttpResponseCode, err error,
) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return httpResponseCode, errors.New("HttpResponseCodeMustBeUint")
	}

	re := regexp.MustCompile(responseCodeExpression)
	stringValue := strconv.FormatUint(uintValue, 10)
	if !re.MatchString(stringValue) {
		return 0, errors.New("InvalidHttpResponseCode")
	}

	return HttpResponseCode(uintValue), nil
}

func NewHttpResponseCodePanic(value interface{}) HttpResponseCode {
	responseCode, err := NewHttpResponseCode(value)
	if err != nil {
		panic(err)
	}

	return responseCode
}

func (vo HttpResponseCode) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
