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

	targetServiceName, err := NewServiceName(mtvStr)
	if err == nil {
		return MappingTargetValue(targetServiceName.String()), nil
	}

	targetUrl, err := NewUrl(mtvStr)
	if err == nil {
		return MappingTargetValue(targetUrl.String()), nil
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

func (mtv MappingTargetValue) String() string {
	return string(mtv)
}
