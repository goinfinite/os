package valueObject

import (
	"errors"
	"strings"
)

type MktplaceItemInstallStep string

func NewMktplaceItemInstallStep(value string) (MktplaceItemInstallStep, error) {
	miis := MktplaceItemInstallStep(value)
	if !miis.isValid() {
		return "", errors.New("InvalidMktItemInstallStep")
	}

	return miis, nil
}

func NewMktplaceItemInstallStepPanic(value string) MktplaceItemInstallStep {
	miis, err := NewMktplaceItemInstallStep(value)
	if err != nil {
		panic(err)
	}

	return miis
}

func (miis MktplaceItemInstallStep) isValid() bool {
	isTooShort := len(string(miis)) < 1
	isTooLong := len(string(miis)) > 512
	return !isTooShort && !isTooLong
}

func (miis MktplaceItemInstallStep) String() string {
	return string(miis)
}

func (miisPtr *MktplaceItemInstallStep) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	miis, err := NewMktplaceItemInstallStep(unquotedValue)
	if err != nil {
		return err
	}

	*miisPtr = miis
	return nil
}
