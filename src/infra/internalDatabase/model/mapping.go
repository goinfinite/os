package dbModel

import (
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type Mapping struct {
	ID                            uint64 `gorm:"primaryKey"`
	Hostname                      string `gorm:"not null"`
	Path                          string `gorm:"not null"`
	MatchPattern                  string `gorm:"not null"`
	TargetType                    string `gorm:"not null"`
	TargetValue                   *string
	TargetHttpResponseCode        *string
	ShouldUpgradeInsecureRequests *bool
	MarketplaceInstalledItemID    *uint
	MarketplaceInstalledItemName  *string
	MappingSecurityRuleID         *uint64
	CreatedAt                     time.Time `gorm:"not null"`
	UpdatedAt                     time.Time `gorm:"not null"`
}

func (Mapping) TableName() string {
	return "mappings"
}

func (model Mapping) ToEntity() (mappingEntity entity.Mapping, err error) {
	mappingId, err := valueObject.NewMappingId(model.ID)
	if err != nil {
		return mappingEntity, err
	}

	hostname, err := valueObject.NewFqdn(model.Hostname)
	if err != nil {
		return mappingEntity, err
	}

	path, err := valueObject.NewMappingPath(model.Path)
	if err != nil {
		return mappingEntity, err
	}

	matchPattern, err := valueObject.NewMappingMatchPattern(model.MatchPattern)
	if err != nil {
		return mappingEntity, err
	}

	targetType, err := valueObject.NewMappingTargetType(model.TargetType)
	if err != nil {
		return mappingEntity, err
	}

	var targetValuePtr *valueObject.MappingTargetValue
	if model.TargetValue != nil && targetType != valueObject.MappingTargetTypeResponseCode {
		targetValue, err := valueObject.NewMappingTargetValue(
			*model.TargetValue, targetType,
		)
		if err != nil {
			return mappingEntity, err
		}
		targetValuePtr = &targetValue
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	if model.TargetHttpResponseCode != nil {
		targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
			*model.TargetHttpResponseCode,
		)
		if err != nil {
			return mappingEntity, err
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var marketplaceInstalledItemIdPtr *valueObject.MarketplaceItemId
	if model.MarketplaceInstalledItemID != nil {
		marketplaceInstalledItemId, err := valueObject.NewMarketplaceItemId(
			*model.MarketplaceInstalledItemID,
		)
		if err != nil {
			return mappingEntity, err
		}
		marketplaceInstalledItemIdPtr = &marketplaceInstalledItemId
	}

	var marketplaceInstalledItemNamePtr *valueObject.MarketplaceItemName
	if model.MarketplaceInstalledItemName != nil {
		marketplaceInstalledItemName, err := valueObject.NewMarketplaceItemName(
			*model.MarketplaceInstalledItemName,
		)
		if err != nil {
			return mappingEntity, err
		}
		marketplaceInstalledItemNamePtr = &marketplaceInstalledItemName
	}

	var mappingSecurityRuleIdPtr *valueObject.MappingSecurityRuleId
	if model.MappingSecurityRuleID != nil {
		mappingSecurityRuleId, err := valueObject.NewMappingSecurityRuleId(
			*model.MappingSecurityRuleID,
		)
		if err != nil {
			return mappingEntity, err
		}
		mappingSecurityRuleIdPtr = &mappingSecurityRuleId
	}

	return entity.NewMapping(
		mappingId, hostname, path, matchPattern, targetType, targetValuePtr,
		targetHttpResponseCodePtr, model.ShouldUpgradeInsecureRequests,
		marketplaceInstalledItemIdPtr, marketplaceInstalledItemNamePtr,
		mappingSecurityRuleIdPtr, valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
