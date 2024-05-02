package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type DataFieldLabel string

var validDataFieldLabels = []string{
	"checkbock",
	"color",
	"date",
	"datetime-local",
	"email",
	"image",
	"month",
	"number",
	"password",
	"radio",
	"range",
	"search",
	"tel",
	"text",
	"time",
	"url",
	"week",
}

func NewDataFieldLabel(value string) (DataFieldLabel, error) {
	value = strings.ToLower(value)
	if !slices.Contains(validDataFieldLabels, value) {
		return "", errors.New("InvalidServiceNature")
	}

	return DataFieldLabel(value), nil
}

func NewDataFieldLabelPanic(value string) DataFieldLabel {
	dfl, err := NewDataFieldLabel(value)
	if err != nil {
		panic(err)
	}

	return dfl
}

func (dfl DataFieldLabel) String() string {
	return string(dfl)
}
