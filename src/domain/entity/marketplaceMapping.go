package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
	"gopkg.in/yaml.v3"
)

type MarketplaceMapping struct {
	Path                    valueObject.MappingPath         `json:"path"`
	MatchPattern            valueObject.MappingMatchPattern `json:"matchPattern"`
	TargetType              valueObject.MappingTargetType   `json:"targetType"`
	TargetServiceName       *valueObject.ServiceName        `json:"targetServiceName,omitempty"`
	TargetUrl               *valueObject.Url                `json:"targetUrl,omitempty"`
	TargetHttpResponseCode  *valueObject.HttpResponseCode   `json:"targetHttpResponseCode,omitempty"`
	TargetInlineHtmlContent *valueObject.InlineHtmlContent  `json:"targetInlineHtmlContent,omitempty"`
}

func NewMarketplaceMapping(
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetServiceName *valueObject.ServiceName,
	targetUrl *valueObject.Url,
	targetHttpResponseCode *valueObject.HttpResponseCode,
	targetInlineHtmlContent *valueObject.InlineHtmlContent,
) MarketplaceMapping {
	return MarketplaceMapping{
		Path:                    path,
		MatchPattern:            matchPattern,
		TargetType:              targetType,
		TargetServiceName:       targetServiceName,
		TargetUrl:               targetUrl,
		TargetHttpResponseCode:  targetHttpResponseCode,
		TargetInlineHtmlContent: targetInlineHtmlContent,
	}
}

func (mmPtr *MarketplaceMapping) UnmarshalYAML(value *yaml.Node) error {
	var valuesMap map[string]string
	err := value.Decode(&valuesMap)
	if err != nil {
		return err
	}

	path, err := valueObject.NewMappingPath(valuesMap["path"])
	if err != nil {
		return err
	}

	matchPattern, err := valueObject.NewMappingMatchPattern(
		valuesMap["matchPattern"],
	)
	if err != nil {
		return err
	}

	targetType, err := valueObject.NewMappingTargetType(
		valuesMap["targetType"],
	)
	if err != nil {
		return err
	}

	var targetSvcNamePtr *valueObject.ServiceName
	targetSvcName, targetSvcNameExists := valuesMap["targetServiceName"]
	if targetSvcNameExists {
		targetSvcName, err := valueObject.NewServiceName(targetSvcName)
		if err != nil {
			return err
		}

		targetSvcNamePtr = &targetSvcName
	}

	var targetUrlPtr *valueObject.Url
	targetUrl, targetUrlExists := valuesMap["targetUrl"]
	if targetUrlExists {
		targetUrl, err := valueObject.NewUrl(targetUrl)
		if err != nil {
			return err
		}

		targetUrlPtr = &targetUrl
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	targetHttpResponseCode, targetHttpResponseCodeExists := valuesMap["targetHttpResponseCode"]
	if targetHttpResponseCodeExists {
		targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
			targetHttpResponseCode,
		)
		if err != nil {
			return err
		}

		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var targetHtmlInlineContentPtr *valueObject.InlineHtmlContent
	targetHtmlInlineContent, targetHtmlInlineContentExists := valuesMap["targetHtmlInlineContent"]
	if targetHtmlInlineContentExists {
		targetHtmlInlineContent, err := valueObject.NewInlineHtmlContent(
			targetHtmlInlineContent,
		)
		if err != nil {
			return err
		}

		targetHtmlInlineContentPtr = &targetHtmlInlineContent
	}

	*mmPtr = NewMarketplaceMapping(
		path,
		matchPattern,
		targetType,
		targetSvcNamePtr,
		targetUrlPtr,
		targetHttpResponseCodePtr,
		targetHtmlInlineContentPtr,
	)

	return nil
}
