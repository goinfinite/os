package valueObject

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
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

func (miisPtr *MarketplaceItemInstallStep) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")
	valueWithoutBackslash := strings.ReplaceAll(unquotedValue, "\\", "")

	miis, err := NewMarketplaceItemInstallStep(valueWithoutBackslash)
	if err != nil {
		return err
	}

	*miisPtr = miis
	return nil
}

func (miisPtr *MarketplaceItemInstallStep) UnmarshalYAML(value *yaml.Node) error {
	var valueStr string
	err := value.Decode(&valueStr)
	if err != nil {
		return err
	}

	miis, err := NewMarketplaceItemInstallStep(valueStr)
	if err != nil {
		return err
	}

	*miisPtr = miis
	return nil
}
