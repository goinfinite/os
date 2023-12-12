package valueObject

import (
	"errors"
	"regexp"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const responseCodeExpression = "^([1-5][0-9][0-9])$"

type HttpResponseCode uint64

func NewHttpResponseCode(value interface{}) (HttpResponseCode, error) {
	responseCodeUint, err := voHelper.InterfaceToUint(value)
	if err != nil {
		return 0, errors.New("InvalidHttpResponseCode")
	}

	responseCode := HttpResponseCode(responseCodeUint)
	if !responseCode.isValid() {
		return 0, errors.New("InvalidHttpResponseCode")
	}

	return responseCode, nil
}

func NewHttpResponseCodePanic(value string) HttpResponseCode {
	responseCode, err := NewHttpResponseCode(value)
	if err != nil {
		panic(err)
	}

	return responseCode
}

func (responseCode HttpResponseCode) isValid() bool {
	responseCodeRegex := regexp.MustCompile(responseCodeExpression)
	return responseCodeRegex.MatchString(responseCode.String())
}

func (responseCode HttpResponseCode) String() string {
	return strconv.FormatUint(uint64(responseCode), 10)
}
