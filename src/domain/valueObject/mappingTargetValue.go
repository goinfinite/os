package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MappingTargetValue string

func NewMappingTargetValue(
	value interface{},
	targetType MappingTargetType,
) (vo MappingTargetValue, err error) {
	voStr, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("InvalidMappingTargetValue")
	}
	voStr = strings.TrimSpace(voStr)

	switch targetType.String() {
	case "url":
		targetUrl, err := NewUrl(voStr)
		if err == nil {
			return MappingTargetValue(targetUrl.String()), nil
		}
	case "service":
		targetServiceName, err := NewServiceName(voStr)
		if err == nil {
			return MappingTargetValue(targetServiceName.String()), nil
		}
	case "response-code":
		targetHttpResponseCode, err := NewHttpResponseCode(voStr)
		if err == nil {
			return MappingTargetValue(
				targetHttpResponseCode.String(),
			), nil
		}
	case "inline-html":
		targetInlineHtmlContent, err := NewInlineHtmlContent(voStr)
		if err == nil {
			return MappingTargetValue(
				targetInlineHtmlContent.String(),
			), nil
		}
	}

	return vo, errors.New("InvalidMappingTargetValue")
}

func NewMappingTargetValuePanic(
	value interface{},
	targetType MappingTargetType,
) (vo MappingTargetValue) {
	vo, err := NewMappingTargetValue(value, targetType)
	if err != nil {
		panic(err)
	}

	return vo
}

func (vo MappingTargetValue) String() string {
	return string(vo)
}
