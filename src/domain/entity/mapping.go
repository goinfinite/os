package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type Mapping struct {
	Id                            valueObject.MappingId              `json:"id"`
	Hostname                      tkValueObject.Fqdn                 `json:"hostname"`
	Path                          valueObject.MappingPath            `json:"path"`
	MatchPattern                  valueObject.MappingMatchPattern    `json:"matchPattern"`
	TargetType                    valueObject.MappingTargetType      `json:"targetType"`
	TargetValue                   *valueObject.MappingTargetValue    `json:"targetValue"`
	TargetHttpResponseCode        *tkValueObject.HttpStatusCode      `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool                              `json:"shouldUpgradeInsecureRequests"`
	MarketplaceInstalledItemId    *valueObject.MarketplaceItemId     `json:"marketplaceInstalledItemId"`
	MarketplaceInstalledItemName  *valueObject.MarketplaceItemName   `json:"marketplaceInstalledItemName"`
	MappingSecurityRuleId         *valueObject.MappingSecurityRuleId `json:"mappingSecurityRuleId"`
	CreatedAt                     tkValueObject.UnixTime             `json:"createdAt"`
	UpdatedAt                     tkValueObject.UnixTime             `json:"updatedAt"`
}

func NewMapping(
	id valueObject.MappingId,
	hostname tkValueObject.Fqdn,
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetValue *valueObject.MappingTargetValue,
	targetHttpResponseCode *tkValueObject.HttpStatusCode,
	shouldUpgradeInsecureRequests *bool,
	marketplaceInstalledItemId *valueObject.MarketplaceItemId,
	marketplaceInstalledItemName *valueObject.MarketplaceItemName,
	mappingSecurityRuleId *valueObject.MappingSecurityRuleId,
	createdAt tkValueObject.UnixTime,
	updatedAt tkValueObject.UnixTime,
) Mapping {
	return Mapping{
		Id:                            id,
		Hostname:                      hostname,
		Path:                          path,
		MatchPattern:                  matchPattern,
		TargetType:                    targetType,
		TargetValue:                   targetValue,
		TargetHttpResponseCode:        targetHttpResponseCode,
		ShouldUpgradeInsecureRequests: shouldUpgradeInsecureRequests,
		MarketplaceInstalledItemId:    marketplaceInstalledItemId,
		MarketplaceInstalledItemName:  marketplaceInstalledItemName,
		MappingSecurityRuleId:         mappingSecurityRuleId,
		CreatedAt:                     createdAt,
		UpdatedAt:                     updatedAt,
	}
}
