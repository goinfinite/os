package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type InlineHtmlContent string

func NewInlineHtmlContent(value interface{}) (
	inlineHtmlContent InlineHtmlContent, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return inlineHtmlContent, errors.New("InlineHtmlContentMustBeString")
	}

	if len(stringValue) == 0 {
		return inlineHtmlContent, errors.New("InlineHtmlContentTooSmall")
	}

	if len(stringValue) > 3500 {
		return inlineHtmlContent, errors.New("InlineHtmlContentTooBig")
	}

	return InlineHtmlContent(stringValue), nil
}

func (vo InlineHtmlContent) String() string {
	return string(vo)
}
