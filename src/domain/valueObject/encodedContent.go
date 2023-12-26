package valueObject

import (
	"encoding/base64"
	"errors"
	"regexp"
)

const encodedContentRegexExpression = `^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}={2})$`

type EncodedContent string

func NewEncodedContent(value string) (EncodedContent, error) {
	isEmpty := len(value) < 1
	if isEmpty {
		return "", errors.New("InvalidEncodedContent")
	}

	encodedBaseContentRegex := regexp.MustCompile(encodedContentRegexExpression)
	isValid := encodedBaseContentRegex.MatchString(value)

	if !isValid {
		return "", errors.New("InvalidEncodedContent")
	}

	return EncodedContent(value), nil
}

func NewEncodedContentPanic(value string) EncodedContent {
	encodedContent, err := NewEncodedContent(value)
	if err != nil {
		panic(err)
	}
	return encodedContent
}

func (encodedBaseContent EncodedContent) GetDecodedContent() (string, error) {
	decodedContent, err := base64.StdEncoding.DecodeString(string(encodedBaseContent))
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}

func (encodedBaseContent EncodedContent) String() string {
	return string(encodedBaseContent)
}
