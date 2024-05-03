package valueObject

import (
	"errors"
	"slices"
	"strings"
)

type HtmlType string

var validHtmlTypes = []string{
	"checkbox",
	"color",
	"date",
	"email",
	"image",
	"number",
	"password",
	"radio",
	"range",
	"search",
	"tel",
	"text",
	"time",
	"url",
}

func NewHtmlType(value string) (HtmlType, error) {
	value = strings.ToLower(value)
	if !slices.Contains(validHtmlTypes, value) {
		return "", errors.New("InvalidHtmlType")
	}

	return HtmlType(value), nil
}

func NewHtmlTypePanic(value string) HtmlType {
	ht, err := NewHtmlType(value)
	if err != nil {
		panic(err)
	}

	return ht
}

func (ht HtmlType) String() string {
	return string(ht)
}
