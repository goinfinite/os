package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type Mapping struct {
	Id                           valueObject.MappingId            `json:"id"`
	Hostname                     valueObject.Fqdn                 `json:"-"`
	Path                         valueObject.MappingPath          `json:"path"`
	MatchPattern                 valueObject.MappingMatchPattern  `json:"matchPattern"`
	TargetType                   valueObject.MappingTargetType    `json:"targetType"`
	TargetValue                  *valueObject.MappingTargetValue  `json:"targetValue"`
	TargetHttpResponseCode       *valueObject.HttpResponseCode    `json:"targetHttpResponseCode"`
	MarketplaceInstalledItemId   *valueObject.MarketplaceItemId   `json:"marketplaceInstalledItemId"`
	MarketplaceInstalledItemName *valueObject.MarketplaceItemName `json:"marketplaceInstalledItemName"`
	CreatedAt                    valueObject.UnixTime             `json:"createdAt"`
	UpdatedAt                    valueObject.UnixTime             `json:"updatedAt"`
}

func NewMapping(
	id valueObject.MappingId,
	hostname valueObject.Fqdn,
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetValue *valueObject.MappingTargetValue,
	targetHttpResponseCode *valueObject.HttpResponseCode,
	marketplaceInstalledItemId *valueObject.MarketplaceItemId,
	marketplaceInstalledItemName *valueObject.MarketplaceItemName,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) Mapping {
	return Mapping{
		Id:                           id,
		Hostname:                     hostname,
		Path:                         path,
		MatchPattern:                 matchPattern,
		TargetType:                   targetType,
		TargetValue:                  targetValue,
		TargetHttpResponseCode:       targetHttpResponseCode,
		MarketplaceInstalledItemId:   marketplaceInstalledItemId,
		MarketplaceInstalledItemName: marketplaceInstalledItemName,
		CreatedAt:                    createdAt,
		UpdatedAt:                    updatedAt,
	}
}
