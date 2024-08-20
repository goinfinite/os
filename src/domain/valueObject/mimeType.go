package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const mimeTypeRegexExpression = `^[A-z0-9\-]{1,64}\/[A-z0-9\-\_\+\.\,]{2,128}$|^(directory|generic)$`

type MimeType string

func NewMimeType(value interface{}) (mimeType MimeType, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return mimeType, errors.New("MimeTypeMustBeString")
	}

	re := regexp.MustCompile(mimeTypeRegexExpression)
	if !re.MatchString(stringValue) {
		return mimeType, errors.New("InvalidMimeType")
	}

	return MimeType(stringValue), nil
}

func (vo MimeType) String() string {
	return string(vo)
}

func (vo MimeType) IsDir() bool {
	return vo == "directory"
}
