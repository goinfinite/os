package valueObject

import (
	"encoding/base64"
	"errors"
	"regexp"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const encodedContentRegexExpression = `^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}={2})$`

type EncodedContent string

func NewEncodedContent(value interface{}) (encodedContent EncodedContent, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return encodedContent, errors.New("EncodedContentMustBeString")
	}

	if len(stringValue) == 0 {
		return encodedContent, errors.New("EncodedContentEmpty")
	}

	re := regexp.MustCompile(encodedContentRegexExpression)
	if !re.MatchString(stringValue) {
		return encodedContent, errors.New("InvalidEncodedContent")
	}

	return EncodedContent(stringValue), nil
}

func (vo EncodedContent) GetDecodedContent() (string, error) {
	decodedContent, err := base64.StdEncoding.DecodeString(string(vo))
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}

func (vo EncodedContent) String() string {
	return string(vo)
}
