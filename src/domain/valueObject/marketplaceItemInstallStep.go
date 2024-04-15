package valueObject

import (
	"errors"
)

type MarketplaceItemCmdStep string

func NewMarketplaceItemCmdStep(value string) (MarketplaceItemCmdStep, error) {
	if len(value) < 1 {
		return "", errors.New("MarketplaceItemCmdStepTooSmall")
	}

	if len(value) > 4096 {
		return "", errors.New("MarketplaceItemCmdStepTooBig")
	}

	return MarketplaceItemCmdStep(value), nil
}

func NewMarketplaceItemCmdStepPanic(value string) MarketplaceItemCmdStep {
	miis, err := NewMarketplaceItemCmdStep(value)
	if err != nil {
		panic(err)
	}

	return miis
}

func (miis MarketplaceItemCmdStep) String() string {
	return string(miis)
}
