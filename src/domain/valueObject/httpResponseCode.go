package valueObject

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

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

func NewHttpResponseCodePanic(value interface{}) HttpResponseCode {
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

func (responseCodePtr *HttpResponseCode) UnmarshalJSON(
	value []byte,
) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	responseCode, err := NewHttpResponseCode(unquotedValue)
	if err != nil {
		return err
	}

	*responseCodePtr = responseCode
	return nil
}
