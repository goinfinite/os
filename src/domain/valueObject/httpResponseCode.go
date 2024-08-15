package valueObject

import (
	"errors"
	"regexp"
	"strconv"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const responseCodeRegex = "^([1-5][0-9][0-9])$"

type HttpResponseCode uint64

func NewHttpResponseCode(value interface{}) (
	httpResponseCode HttpResponseCode, err error,
) {
	uintValue, err := voHelper.InterfaceToUint64(value)
	if err != nil {
		return httpResponseCode, errors.New("HttpResponseCodeMustBeUint")
	}
	stringValue := strconv.FormatUint(uintValue, 10)

	re := regexp.MustCompile(responseCodeRegex)
	if !re.MatchString(stringValue) {
		return httpResponseCode, errors.New("InvalidHttpResponseCode")
	}

	return HttpResponseCode(uintValue), nil
}

func (vo HttpResponseCode) String() string {
	return strconv.FormatUint(uint64(vo), 10)
}
