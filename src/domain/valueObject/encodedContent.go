package valueObject

import (
	"encoding/base64"
	"errors"
	"regexp"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var encodedContentRegex = regexp.MustCompile(`^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}={2})$`)

type EncodedContent string

func NewEncodedContent(value any) (encodedContent EncodedContent, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return encodedContent, errors.New("EncodedContentMustBeString")
	}

	if len(stringValue) == 0 {
		return encodedContent, errors.New("EmptyEncodedContent")
	}

	if !encodedContentRegex.MatchString(stringValue) {
		return encodedContent, errors.New("InvalidEncodedContent")
	}

	return EncodedContent(stringValue), nil
}

func (vo EncodedContent) DecodeContent() (voStr string, err error) {
	decodedContent, err := base64.StdEncoding.DecodeString(string(vo))
	if err != nil {
		return voStr, err
	}

	return string(decodedContent), nil
}

func (vo EncodedContent) String() string {
	return string(vo)
}
