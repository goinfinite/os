package entity

import "github.com/speedianet/os/src/domain/valueObject"

type PhpVersion struct {
	Value   valueObject.PhpVersion   `json:"value"`
	Options []valueObject.PhpVersion `json:"options"`
}

func NewPhpVersion(
	value valueObject.PhpVersion,
	options []valueObject.PhpVersion,
) PhpVersion {
	return PhpVersion{
		Value:   value,
		Options: options,
	}
}
