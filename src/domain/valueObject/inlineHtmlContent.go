package valueObject

import (
	"errors"
)

type InlineHtmlContent string

func NewInlineHtmlContent(value string) (InlineHtmlContent, error) {
	if len(value) == 0 || len(value) > 3000 {
		return "", errors.New("InvalidInlineHtmlContent")
	}

	return InlineHtmlContent(value), nil
}

func NewInlineHtmlContentPanic(value string) InlineHtmlContent {
	inlineHtmlContent, err := NewInlineHtmlContent(value)
	if err != nil {
		panic(err)
	}

	return inlineHtmlContent
}

func (ihc InlineHtmlContent) String() string {
	return string(ihc)
}
