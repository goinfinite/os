package valueObject

import (
	"errors"
	"slices"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type MappingTargetType string

const (
	MappingTargetTypeUrl          MappingTargetType = "url"
	MappingTargetTypeService      MappingTargetType = "service"
	MappingTargetTypeResponseCode MappingTargetType = "response-code"
	MappingTargetTypeInlineHtml   MappingTargetType = "inline-html"
	MappingTargetTypeStaticFiles  MappingTargetType = "static-files"
)

var ValidMappingTargetTypes = []string{
	MappingTargetTypeUrl.String(), MappingTargetTypeService.String(),
	MappingTargetTypeResponseCode.String(), MappingTargetTypeInlineHtml.String(),
	MappingTargetTypeStaticFiles.String(),
}

func NewMappingTargetType(value interface{}) (
	mappingTargetType MappingTargetType, err error,
) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return mappingTargetType, errors.New("MappingTargetTypeMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidMappingTargetTypes, stringValue) {
		return mappingTargetType, errors.New("InvalidMappingTargetType")
	}

	return MappingTargetType(stringValue), nil
}

func (vo MappingTargetType) String() string {
	return string(vo)
}
