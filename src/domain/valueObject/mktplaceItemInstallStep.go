package valueObject

import "errors"

type MktplaceItemInstallStep string

func NewMktplaceItemInstallStep(value string) (MktplaceItemInstallStep, error) {
	mkplaceItemInstallStep := MktplaceItemInstallStep(value)
	if !mkplaceItemInstallStep.isValid() {
		return "", errors.New("InvalidMarketplaceItemInstallStep")
	}

	return mkplaceItemInstallStep, nil
}

func NewMktplaceItemInstallStepPanic(value string) MktplaceItemInstallStep {
	mkplaceItemInstallStep, err := NewMktplaceItemInstallStep(value)
	if err != nil {
		panic(err)
	}

	return mkplaceItemInstallStep
}

func (mkplaceItemInstallStep MktplaceItemInstallStep) isValid() bool {
	isTooShort := len(string(mkplaceItemInstallStep)) < 1
	isTooLong := len(string(mkplaceItemInstallStep)) > 512
	return !isTooShort && !isTooLong
}

func (mkplaceItemInstallStep MktplaceItemInstallStep) String() string {
	return string(mkplaceItemInstallStep)
}
