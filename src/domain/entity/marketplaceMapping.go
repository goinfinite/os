package entity

import "github.com/speedianet/os/src/domain/valueObject"

type MarketplaceMapping struct {
	Path              valueObject.MappingPath         `json:"path"`
	MatchPattern      valueObject.MappingMatchPattern `json:"matchPattern"`
	TargetServiceName *valueObject.ServiceName        `json:"targetServiceName"`
}

func NewMarketplaceMapping(
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetServiceName *valueObject.ServiceName,
) MarketplaceMapping {
	return MarketplaceMapping{
		Path:              path,
		MatchPattern:      matchPattern,
		TargetServiceName: targetServiceName,
	}
}
