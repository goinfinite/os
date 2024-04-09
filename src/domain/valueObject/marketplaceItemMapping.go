package valueObject

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
