package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type InlineHtmlContent string

func NewInlineHtmlContent(value interface{}) (
	inlineHtmlContent InlineHtmlContent, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return inlineHtmlContent, errors.New("InlineHtmlContentMustBeString")
	}

	if len(stringValue) == 0 {
		return inlineHtmlContent, errors.New("EmptyInlineHtmlContent")
	}

	if len(stringValue) > 3500 {
		return inlineHtmlContent, errors.New("InlineHtmlContentTooBig")
	}

	return InlineHtmlContent(stringValue), nil
}

func (vo InlineHtmlContent) String() string {
	return string(vo)
}
