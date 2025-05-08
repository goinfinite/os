package valueObject

type MarketplaceItemMapping struct {
	Path                          MappingPath         `json:"path"`
	MatchPattern                  MappingMatchPattern `json:"matchPattern"`
	TargetType                    MappingTargetType   `json:"targetType"`
	TargetValue                   *MappingTargetValue `json:"targetValue"`
	TargetHttpResponseCode        *HttpResponseCode   `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool               `json:"shouldUpgradeInsecureRequests"`
}

func NewMarketplaceItemMapping(
	path MappingPath,
	matchPattern MappingMatchPattern,
	targetType MappingTargetType,
	targetValue *MappingTargetValue,
	targetHttpResponseCode *HttpResponseCode,
	shouldUpgradeInsecureRequests *bool,
) MarketplaceItemMapping {
	return MarketplaceItemMapping{
		Path:                          path,
		MatchPattern:                  matchPattern,
		TargetType:                    targetType,
		TargetValue:                   targetValue,
		TargetHttpResponseCode:        targetHttpResponseCode,
		ShouldUpgradeInsecureRequests: shouldUpgradeInsecureRequests,
	}
}
