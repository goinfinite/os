package valueObject

import (
	"encoding/base64"
	"errors"
	"regexp"
)

const encodedContentRegexExpression = `^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}={2})$`

type EncodedContent string

func NewEncodedContent(value string) (EncodedContent, error) {
	encodedBaseContent := EncodedContent(value)
	if !encodedBaseContent.isValid() {
		return "", errors.New("InvalidEncodedContent")
	}
	return encodedBaseContent, nil
}

func NewEncodedContentPanic(value string) EncodedContent {
	encodedBaseContent, err := NewEncodedContent(value)
	if err != nil {
		panic(err)
	}
	return encodedBaseContent
}

func (encodedBaseContent EncodedContent) isValid() bool {
	isEmpty := false
	if len(encodedBaseContent) < 1 {
		isEmpty = true
	}

	encodedBaseContentRegex := regexp.MustCompile(encodedContentRegexExpression)
	return encodedBaseContentRegex.MatchString(string(encodedBaseContent)) && !isEmpty
}

func (encodedBaseContent EncodedContent) GetDecodedContent() string {
	decodedContent, err := base64.StdEncoding.DecodeString(string(encodedBaseContent))
	if err != nil {
		return ""
	}

	return string(decodedContent)
}

func (encodedBaseContent EncodedContent) String() string {
	return string(encodedBaseContent)
}
