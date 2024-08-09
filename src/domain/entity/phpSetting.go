package entity

import (
	"errors"
	"strings"

	"github.com/speedianet/os/src/domain/valueObject"
)

type PhpSetting struct {
	Name    valueObject.PhpSettingName     `json:"name"`
	Value   valueObject.PhpSettingValue    `json:"value"`
	Options []valueObject.PhpSettingOption `json:"options"`
}

func NewPhpSetting(
	name valueObject.PhpSettingName,
	value valueObject.PhpSettingValue,
	options []valueObject.PhpSettingOption,
) PhpSetting {
	return PhpSetting{
		Name:    name,
		Value:   value,
		Options: options,
	}
}

// format: name:value:suggestedValue1,suggestedValue2,suggestedValue3
func NewPhpSettingFromString(stringValue string) (setting PhpSetting, err error) {
	stringValueParts := strings.Split(stringValue, ":")
	if len(stringValueParts) == 0 {
		return setting, errors.New("PhpSettingEmpty")
	}

	if len(stringValueParts) < 2 {
		return setting, errors.New("MissingPhpSettingParts")
	}

	name, err := valueObject.NewPhpSettingName(stringValueParts[0])
	if err != nil {
		return setting, err
	}

	value, err := valueObject.NewPhpSettingValue(stringValueParts[1])
	if err != nil {
		return setting, err
	}

	options := []valueObject.PhpSettingOption{}

	if len(stringValueParts) == 2 {
		return NewPhpSetting(name, value, options), nil
	}

	optionsParts := strings.Split(stringValueParts[2], ",")
	if len(optionsParts) > 0 {
		for _, optionStr := range optionsParts {
			option, err := valueObject.NewPhpSettingOption(optionStr)
			if err != nil {
				continue
			}
			options = append(options, option)
		}
	}

	return NewPhpSetting(name, value, options), nil
}
