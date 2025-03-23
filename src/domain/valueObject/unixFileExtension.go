package valueObject

import (
	"errors"
	"mime"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const unixFileExtensionRegexExpression = `^[\w\_\-]{1,15}$`

type UnixFileExtension string

func NewUnixFileExtension(value interface{}) (
	unixFileExtension UnixFileExtension, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return unixFileExtension, errors.New("UnixFileExtensionMustBeString")
	}

	if strings.HasPrefix(stringValue, ".") {
		stringValue, _ = strings.CutPrefix(stringValue, ".")
	}

	re := regexp.MustCompile(unixFileExtensionRegexExpression)
	if !re.MatchString(stringValue) {
		return unixFileExtension, errors.New("InvalidUnixFileExtension")
	}

	return UnixFileExtension(stringValue), nil
}

func (vo UnixFileExtension) GetMimeType() MimeType {
	mimeTypeStr := MimeTypeGeneric.String()

	fileExtWithLeadingDot := "." + string(vo)
	mimeTypeWithCharset := mime.TypeByExtension(fileExtWithLeadingDot)
	if len(mimeTypeWithCharset) > 1 {
		mimeTypeOnly := strings.Split(mimeTypeWithCharset, ";")[0]
		mimeTypeStr = mimeTypeOnly
	}

	mimeType, _ := NewMimeType(mimeTypeStr)
	return mimeType
}

func (vo UnixFileExtension) String() string {
	return string(vo)
}
