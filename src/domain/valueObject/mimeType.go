package valueObject

import (
	"errors"
	"regexp"
)

const mimeTypeRegexExpression = `^[\p{L}0-9\-]{1,64}\/[\p{L}0-9\-\_\+\.\,]{2,128}$|^(directory|generic)$`

type MimeType string

func NewMimeType(mimeTypeStr string) (MimeType, error) {
	mimeType := MimeType(mimeTypeStr)
	if !mimeType.isValid() {
		return "", errors.New("InvalidMimeType")
	}
	return mimeType, nil
}

func NewMimeTypePanic(mimeTypeStr string) MimeType {
	mimeType, err := NewMimeType(mimeTypeStr)
	if err != nil {
		panic(err)
	}
	return mimeType
}

func (mimeType MimeType) isValid() bool {
	mimeTypeRegexRegex := regexp.MustCompile(mimeTypeRegexExpression)
	return mimeTypeRegexRegex.MatchString(string(mimeType))
}

func (mimeType MimeType) String() string {
	return string(mimeType)
}
