package valueObject

import (
	"gopkg.in/yaml.v3"
)

type MarketplaceItemMapping struct {
	Path                    MappingPath         `json:"path"`
	MatchPattern            MappingMatchPattern `json:"matchPattern"`
	TargetType              MappingTargetType   `json:"targetType"`
	TargetServiceName       *ServiceName        `json:"targetServiceName,omitempty"`
	TargetUrl               *Url                `json:"targetUrl,omitempty"`
	TargetHttpResponseCode  *HttpResponseCode   `json:"targetHttpResponseCode,omitempty"`
	TargetInlineHtmlContent *InlineHtmlContent  `json:"targetInlineHtmlContent,omitempty"`
}

func NewMarketplaceItemMapping(
	path MappingPath,
	matchPattern MappingMatchPattern,
	targetType MappingTargetType,
	targetServiceName *ServiceName,
	targetUrl *Url,
	targetHttpResponseCode *HttpResponseCode,
	targetInlineHtmlContent *InlineHtmlContent,
) MarketplaceItemMapping {
	return MarketplaceItemMapping{
		Path:                    path,
		MatchPattern:            matchPattern,
		TargetType:              targetType,
		TargetServiceName:       targetServiceName,
		TargetUrl:               targetUrl,
		TargetHttpResponseCode:  targetHttpResponseCode,
		TargetInlineHtmlContent: targetInlineHtmlContent,
	}
}

func (mmPtr *MarketplaceItemMapping) UnmarshalYAML(value *yaml.Node) error {
	var valuesMap map[string]string
	err := value.Decode(&valuesMap)
	if err != nil {
		return err
	}

	path, err := NewMappingPath(valuesMap["path"])
	if err != nil {
		return err
	}

	matchPattern, err := NewMappingMatchPattern(
		valuesMap["matchPattern"],
	)
	if err != nil {
		return err
	}

	targetType, err := NewMappingTargetType(
		valuesMap["targetType"],
	)
	if err != nil {
		return err
	}

	var targetSvcNamePtr *ServiceName
	targetSvcName, targetSvcNameExists := valuesMap["targetServiceName"]
	if targetSvcNameExists {
		targetSvcName, err := NewServiceName(targetSvcName)
		if err != nil {
			return err
		}

		targetSvcNamePtr = &targetSvcName
	}

	var targetUrlPtr *Url
	targetUrl, targetUrlExists := valuesMap["targetUrl"]
	if targetUrlExists {
		targetUrl, err := NewUrl(targetUrl)
		if err != nil {
			return err
		}

		targetUrlPtr = &targetUrl
	}

	var targetHttpResponseCodePtr *HttpResponseCode
	targetHttpResponseCode, targetHttpResponseCodeExists := valuesMap["targetHttpResponseCode"]
	if targetHttpResponseCodeExists {
		targetHttpResponseCode, err := NewHttpResponseCode(
			targetHttpResponseCode,
		)
		if err != nil {
			return err
		}

		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var targetHtmlInlineContentPtr *InlineHtmlContent
	targetHtmlInlineContent, targetHtmlInlineContentExists := valuesMap["targetHtmlInlineContent"]
	if targetHtmlInlineContentExists {
		targetHtmlInlineContent, err := NewInlineHtmlContent(
			targetHtmlInlineContent,
		)
		if err != nil {
			return err
		}

		targetHtmlInlineContentPtr = &targetHtmlInlineContent
	}

	*mmPtr = NewMarketplaceItemMapping(
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
