package entity

import (
	"errors"
	"strings"

	"github.com/goinfinite/os/src/domain/valueObject"
)

type PhpSetting struct {
	Name    valueObject.PhpSettingName     `json:"name"`
	Type    valueObject.PhpSettingType     `json:"type"`
	Value   valueObject.PhpSettingValue    `json:"value"`
	Options []valueObject.PhpSettingOption `json:"options"`
}

func NewPhpSetting(
	name valueObject.PhpSettingName,
	settingType valueObject.PhpSettingType,
	value valueObject.PhpSettingValue,
	options []valueObject.PhpSettingOption,
) PhpSetting {
	return PhpSetting{
		Name:    name,
		Type:    settingType,
		Value:   value,
		Options: options,
	}
}

// format: name:value:type:suggestedValue1,suggestedValue2,suggestedValue3
func NewPhpSettingFromString(stringValue string) (setting PhpSetting, err error) {
	stringValueParts := strings.Split(stringValue, ":")
	if len(stringValueParts) == 0 {
		return setting, errors.New("EmptyPhpSetting")
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

	settingType, _ := valueObject.NewPhpSettingType("text")
	options := []valueObject.PhpSettingOption{}

	if len(stringValueParts) == 2 {
		return NewPhpSetting(name, settingType, value, options), nil
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

	settingType, _ = valueObject.NewPhpSettingType("select")
	return NewPhpSetting(name, settingType, value, options), nil
}
