package valueObject

import (
	"errors"
	"regexp"
)

const inlineHtmlContentRegex = `^<html>.*</html>$`

type InlineHtmlContent string

func NewInlineHtmlContent(value string) (InlineHtmlContent, error) {
	inlineHtmlContent := InlineHtmlContent(value)
	if !inlineHtmlContent.isValid() {
		return "", errors.New("InvalidInlineHtmlContent")
	}

	return inlineHtmlContent, nil
}

func NewInlineHtmlContentPanic(value string) InlineHtmlContent {
	inlineHtmlContent, err := NewInlineHtmlContent(value)
	if err != nil {
		panic(err)
	}

	return inlineHtmlContent
}

func (ihc InlineHtmlContent) isValid() bool {
	compiledRegex := regexp.MustCompile(inlineHtmlContentRegex)
	return compiledRegex.MatchString(string(ihc))
}

func (ihc InlineHtmlContent) String() string {
	return string(ihc)
}
