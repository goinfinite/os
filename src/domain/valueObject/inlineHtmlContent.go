package valueObject

import (
	"errors"
	"strings"
)

type InlineHtmlContent string

func NewInlineHtmlContent(value string) (InlineHtmlContent, error) {
	if len(value) == 0 {
		return "", errors.New("InlineHtmlContentTooSmall")
	}

	if len(value) > 3500 {
		return "", errors.New("InlineHtmlContentTooBig")
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

func (ihcPtr *InlineHtmlContent) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	ihc, err := NewInlineHtmlContent(unquotedValue)
	if err != nil {
		return err
	}

	*ihcPtr = ihc
	return nil
}
