package valueObject

import (
	"encoding/base64"
	"errors"
	"regexp"
)

const encodedBase64ContentRegexExpression = `^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}={2})$`

type EncodedBase64Content string

func NewEncodedBase64Content(value string) (EncodedBase64Content, error) {
	encodedBase64 := EncodedBase64Content(value)
	if !encodedBase64.isValid() {
		return "", errors.New("InvalidEncodedBase64")
	}
	return encodedBase64, nil
}

func NewEncodedBase64ContentPanic(value string) EncodedBase64Content {
	EncodedBase64, err := NewEncodedBase64Content(value)
	if err != nil {
		panic(err)
	}
	return EncodedBase64
}

func (encodedBase64 EncodedBase64Content) isValid() bool {
	isEmpty := false
	if len(encodedBase64) < 1 {
		isEmpty = true
	}

	encodedBase64Regex := regexp.MustCompile(encodedBase64ContentRegexExpression)
	return encodedBase64Regex.MatchString(string(encodedBase64)) && !isEmpty
}

func (encodedBase64 EncodedBase64Content) GetDecodedContent() string {
	decodedContent, err := base64.StdEncoding.DecodeString(string(encodedBase64))
	if err != nil {
		return ""
	}

	return string(decodedContent)
}

func (encodedBase64 EncodedBase64Content) String() string {
	return string(encodedBase64)
}
