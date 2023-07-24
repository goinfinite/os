package entity

import "github.com/speedianet/sam/src/domain/valueObject"

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
