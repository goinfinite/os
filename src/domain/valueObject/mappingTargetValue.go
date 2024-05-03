package valueObject

import (
	"errors"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type MappingTargetValue string

func NewMappingTargetValue(value interface{}) (MappingTargetValue, error) {
	mtvStr, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("InvalidMappingTargetValue")
	}

	targetUrl, err := NewUrl(mtvStr)
	if err == nil {
		return MappingTargetValue(targetUrl.String()), nil
	}

	targetServiceName, err := NewServiceName(mtvStr)
	if err == nil {
		return MappingTargetValue(targetServiceName.String()), nil
	}

	targetHttpResponseCode, err := NewHttpResponseCode(mtvStr)
	if err == nil {
		return MappingTargetValue(
			targetHttpResponseCode.String(),
		), nil
	}

	targetInlineHtmlContent, err := NewInlineHtmlContent(mtvStr)
	if err == nil {
		return MappingTargetValue(
			targetInlineHtmlContent.String(),
		), nil
	}

	return "", errors.New("InvalidMappingTargetValue")
}

func NewMappingTargetValuePanic(value interface{}) MappingTargetValue {
	mtv, err := NewMappingTargetValue(value)
	if err != nil {
		panic(err)
	}

	return mtv
}

func NewMappingTargetValueBasedOnType(
	value interface{},
	targetType MappingTargetType,
) (MappingTargetValue, error) {
	var mtv MappingTargetValue

	mtvStr, err := voHelper.InterfaceToString(value)
	if err != nil {
		return "", errors.New("InvalidMappingTargetValue")
	}

	switch targetType.String() {
	case "url":
		targetUrl, err := NewUrl(mtvStr)
		if err == nil {
			return MappingTargetValue(targetUrl.String()), nil
		}
	case "service":
		targetServiceName, err := NewServiceName(mtvStr)
		if err == nil {
			return MappingTargetValue(targetServiceName.String()), nil
		}
	case "response-code":
		targetHttpResponseCode, err := NewHttpResponseCode(mtvStr)
		if err == nil {
			return MappingTargetValue(
				targetHttpResponseCode.String(),
			), nil
		}
	case "inline-html":
		targetInlineHtmlContent, err := NewInlineHtmlContent(mtvStr)
		if err == nil {
			return MappingTargetValue(
				targetInlineHtmlContent.String(),
			), nil
		}
	}

	return mtv, errors.New("InvalidMappingTargetValue")
}

func (mtv MappingTargetValue) String() string {
	return string(mtv)
}
