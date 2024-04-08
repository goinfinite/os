package valueObject

import (
	"errors"
)

type MarketplaceItemInstallStep string

func NewMarketplaceItemInstallStep(value string) (MarketplaceItemInstallStep, error) {
	if len(value) < 1 {
		return "", errors.New("MarketplaceItemInstallStepTooSmall")
	}

	if len(value) > 512 {
		return "", errors.New("MarketplaceItemInstallStepTooBig")
	}

	return MarketplaceItemInstallStep(value), nil
}

func NewMarketplaceItemInstallStepPanic(value string) MarketplaceItemInstallStep {
	miis, err := NewMarketplaceItemInstallStep(value)
	if err != nil {
		panic(err)
	}

	return miis
}

func (miis MarketplaceItemInstallStep) String() string {
	return string(miis)
}
