package valueObject

import (
	"errors"
	"strings"
)

type MktplaceItemInstallStep string

func NewMktplaceItemInstallStep(value string) (MktplaceItemInstallStep, error) {
	mpis := MktplaceItemInstallStep(value)
	if !mpis.isValid() {
		return "", errors.New("InvalidMarketplaceItemInstallStep")
	}

	return mpis, nil
}

func NewMktplaceItemInstallStepPanic(value string) MktplaceItemInstallStep {
	mpis, err := NewMktplaceItemInstallStep(value)
	if err != nil {
		panic(err)
	}

	return mpis
}

func (mpis MktplaceItemInstallStep) isValid() bool {
	isTooShort := len(string(mpis)) < 1
	isTooLong := len(string(mpis)) > 512
	return !isTooShort && !isTooLong
}

func (mpis MktplaceItemInstallStep) String() string {
	return string(mpis)
}

func (mpisPtr *MktplaceItemInstallStep) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mpis, err := NewMktplaceItemInstallStep(unquotedValue)
	if err != nil {
		return err
	}

	*mpisPtr = mpis
	return nil
}
