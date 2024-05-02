package valueObject

import (
	"errors"
	"regexp"
)

const dataFieldLabelRegex string = `^\w[\w-\s]{1,128}\w$`

type DataFieldLabel string

func NewDataFieldLabel(value string) (DataFieldLabel, error) {
	dfl := DataFieldLabel(value)
	if !dfl.isValid() {
		return "", errors.New("InvalidDataFieldLabel")
	}

	return dfl, nil
}

func NewDataFieldLabelPanic(value string) DataFieldLabel {
	dfl, err := NewDataFieldLabel(value)
	if err != nil {
		panic(err)
	}

	return dfl
}

func (dfl DataFieldLabel) isValid() bool {
	re := regexp.MustCompile(dataFieldLabelRegex)
	return re.MatchString(string(dfl))
}

func (dfl DataFieldLabel) String() string {
	return string(dfl)
}
