package valueObject

import (
	"errors"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type MappingTargetValue string

func NewMappingTargetValue(
	value interface{}, targetType MappingTargetType,
) (mappingTargetValue MappingTargetValue, err error) {
	switch targetType.String() {
	case "url":
		targetUrl, err := tkValueObject.NewUrl(value)
		if err == nil {
			return MappingTargetValue(targetUrl.String()), nil
		}
	case "service":
		targetServiceName, err := NewServiceName(value)
		if err == nil {
			return MappingTargetValue(targetServiceName.String()), nil
		}
	case "response-code":
		targetHttpResponseCode, err := tkValueObject.NewHttpStatusCode(value)
		if err == nil {
			return MappingTargetValue(
				targetHttpResponseCode.String(),
			), nil
		}
	case "inline-html":
		targetInlineHtmlContent, err := NewInlineHtmlContent(value)
		if err == nil {
			return MappingTargetValue(
				targetInlineHtmlContent.String(),
			), nil
		}
	}

	return mappingTargetValue, errors.New("InvalidMappingTargetValue")
}

func (vo MappingTargetValue) String() string {
	return string(vo)
}
