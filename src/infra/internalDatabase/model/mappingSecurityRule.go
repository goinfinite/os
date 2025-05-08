package dbModel

import (
	"log/slog"
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type MappingSecurityRule struct {
	ID                             uint64 `gorm:"primaryKey"`
	Name                           string `gorm:"not null"`
	Description                    *string
	AllowedIps                     []string `gorm:"serializer:json"`
	BlockedIps                     []string `gorm:"serializer:json"`
	RpsSoftLimitPerIp              *uint
	RpsHardLimitPerIp              *uint
	ResponseCodeOnMaxRequests      *uint
	MaxConnectionsPerIp            *uint
	BandwidthBpsLimitPerConnection *uint64
	BandwidthLimitOnlyAfterBytes   *uint64
	ResponseCodeOnMaxConnections   *uint
	CreatedAt                      time.Time `gorm:"not null"`
	UpdatedAt                      time.Time `gorm:"not null"`
}

func (MappingSecurityRule) TableName() string {
	return "mapping_security_rules"
}

func (MappingSecurityRule) InitialEntries() (initialEntries []interface{}, err error) {
	for _, initialPreset := range entity.MappingSecurityRuleInitialPresets() {
		initialEntries = append(initialEntries, MappingSecurityRule{}.ToModel(initialPreset))
	}

	return initialEntries, nil
}

func (MappingSecurityRule) ToModel(ruleEntity entity.MappingSecurityRule) MappingSecurityRule {
	var descriptionPtr *string
	if ruleEntity.Description != nil {
		descriptionStr := ruleEntity.Description.String()
		descriptionPtr = &descriptionStr
	}

	allowedIps := []string{}
	for _, ipAddress := range ruleEntity.AllowedIps {
		allowedIps = append(allowedIps, ipAddress.String())
	}

	blockedIps := []string{}
	for _, ipAddress := range ruleEntity.BlockedIps {
		blockedIps = append(blockedIps, ipAddress.String())
	}

	var bandwidthBpsLimitPerConnectionPtr *uint64
	if ruleEntity.BandwidthBpsLimitPerConnection != nil {
		perConnectionUint64 := ruleEntity.BandwidthBpsLimitPerConnection.Uint64()
		bandwidthBpsLimitPerConnectionPtr = &perConnectionUint64
	}

	var bandwidthLimitOnlyAfterBytesPtr *uint64
	if ruleEntity.BandwidthLimitOnlyAfterBytes != nil {
		afterBytesUint64 := ruleEntity.BandwidthLimitOnlyAfterBytes.Uint64()
		bandwidthLimitOnlyAfterBytesPtr = &afterBytesUint64
	}

	return MappingSecurityRule{
		Name:                           ruleEntity.Name.String(),
		Description:                    descriptionPtr,
		AllowedIps:                     allowedIps,
		BlockedIps:                     blockedIps,
		RpsSoftLimitPerIp:              ruleEntity.RpsSoftLimitPerIp,
		RpsHardLimitPerIp:              ruleEntity.RpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      ruleEntity.ResponseCodeOnMaxRequests,
		MaxConnectionsPerIp:            ruleEntity.MaxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: bandwidthBpsLimitPerConnectionPtr,
		BandwidthLimitOnlyAfterBytes:   bandwidthLimitOnlyAfterBytesPtr,
		ResponseCodeOnMaxConnections:   ruleEntity.ResponseCodeOnMaxConnections,
	}
}

func (model MappingSecurityRule) ToEntity() (ruleEntity entity.MappingSecurityRule, err error) {
	id, err := valueObject.NewMappingSecurityRuleId(model.ID)
	if err != nil {
		return ruleEntity, err
	}

	name, err := valueObject.NewMappingSecurityRuleName(model.Name)
	if err != nil {
		return ruleEntity, err
	}

	var descriptionPtr *valueObject.MappingSecurityRuleDescription
	if model.Description != nil {
		description, err := valueObject.NewMappingSecurityRuleDescription(*model.Description)
		if err != nil {
			return ruleEntity, err
		}
		descriptionPtr = &description
	}

	allowedIps := []tkValueObject.CidrBlock{}
	for _, rawIpAddress := range model.AllowedIps {
		ipAddress, err := tkValueObject.NewCidrBlock(rawIpAddress)
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawIpAddress", rawIpAddress))
			continue
		}
		allowedIps = append(allowedIps, ipAddress)
	}

	blockedIps := []tkValueObject.CidrBlock{}
	for _, rawIpAddress := range model.BlockedIps {
		ipAddress, err := tkValueObject.NewCidrBlock(rawIpAddress)
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawIpAddress", rawIpAddress))
			continue
		}
		blockedIps = append(blockedIps, ipAddress)
	}

	var bandwidthBpsLimitPerConnectionPtr *valueObject.Byte
	if model.BandwidthBpsLimitPerConnection != nil {
		bandwidthBpsLimit, err := valueObject.NewByte(*model.BandwidthBpsLimitPerConnection)
		if err != nil {
			return ruleEntity, err
		}
		bandwidthBpsLimitPerConnectionPtr = &bandwidthBpsLimit
	}

	var bandwidthLimitOnlyAfterBytesPtr *valueObject.Byte
	if model.BandwidthLimitOnlyAfterBytes != nil {
		afterBytes, err := valueObject.NewByte(*model.BandwidthLimitOnlyAfterBytes)
		if err != nil {
			return ruleEntity, err
		}
		bandwidthLimitOnlyAfterBytesPtr = &afterBytes
	}

	return entity.NewMappingSecurityRule(
		id, name, descriptionPtr, allowedIps, blockedIps, model.RpsSoftLimitPerIp,
		model.RpsHardLimitPerIp, model.ResponseCodeOnMaxRequests, model.MaxConnectionsPerIp,
		bandwidthBpsLimitPerConnectionPtr, bandwidthLimitOnlyAfterBytesPtr,
		model.ResponseCodeOnMaxConnections, valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
