package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

var (
	MappingSecurityRuleDefaultResponseCodeOnMaxRequests    uint = 429
	MappingSecurityRuleDefaultResponseCodeOnMaxConnections uint = 420
)

type MappingSecurityRule struct {
	Id                             valueObject.MappingSecurityRuleId           `json:"id"`
	Name                           valueObject.MappingSecurityRuleName         `json:"name"`
	Description                    *valueObject.MappingSecurityRuleDescription `json:"description"`
	AllowedIps                     []tkValueObject.CidrBlock                   `json:"allowedIps"`
	BlockedIps                     []tkValueObject.CidrBlock                   `json:"blockedIps"`
	RpsSoftLimitPerIp              *uint                                       `json:"rpsSoftLimitPerIp"`
	RpsHardLimitPerIp              *uint                                       `json:"rpsHardLimitPerIp"`
	ResponseCodeOnMaxRequests      *uint                                       `json:"responseCodeOnMaxRequests"`
	MaxConnectionsPerIp            *uint                                       `json:"maxConnectionsPerIp"`
	BandwidthBpsLimitPerConnection *valueObject.Byte                           `json:"bandwidthBpsLimitPerConnection"`
	BandwidthLimitOnlyAfterBytes   *valueObject.Byte                           `json:"bandwidthLimitOnlyAfterBytes"`
	ResponseCodeOnMaxConnections   *uint                                       `json:"responseCodeOnMaxConnections"`
	CreatedAt                      valueObject.UnixTime                        `json:"createdAt"`
	UpdatedAt                      valueObject.UnixTime                        `json:"updatedAt"`
}

func NewMappingSecurityRule(
	id valueObject.MappingSecurityRuleId,
	name valueObject.MappingSecurityRuleName,
	description *valueObject.MappingSecurityRuleDescription,
	allowedIps, blockedIps []tkValueObject.CidrBlock,
	rpsSoftLimitPerIp, rpsHardLimitPerIp, responseCodeOnMaxRequests, maxConnectionsPerIp *uint,
	bandwidthBpsLimitPerConnection, bandwidthLimitOnlyAfterBytes *valueObject.Byte,
	responseCodeOnMaxConnections *uint,
	createdAt, updatedAt valueObject.UnixTime,
) MappingSecurityRule {
	return MappingSecurityRule{
		Id:                             id,
		Name:                           name,
		Description:                    description,
		AllowedIps:                     allowedIps,
		BlockedIps:                     blockedIps,
		RpsSoftLimitPerIp:              rpsSoftLimitPerIp,
		RpsHardLimitPerIp:              rpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      responseCodeOnMaxRequests,
		MaxConnectionsPerIp:            maxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: bandwidthBpsLimitPerConnection,
		BandwidthLimitOnlyAfterBytes:   bandwidthLimitOnlyAfterBytes,
		ResponseCodeOnMaxConnections:   responseCodeOnMaxConnections,
		CreatedAt:                      createdAt,
		UpdatedAt:                      updatedAt,
	}
}

func MappingSecurityRulePresetBreezy() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Security in shorts and flip‑flops. Best for debugging or very low‑risk scenarios.",
	)
	rpsSoftLimitPerIp := uint(24)
	rpsHardLimitPerIp := uint(32)
	maxConnectionsPerIp := uint(16)
	bandwidthBpsLimitPerConnection := valueObject.Byte(32 * 1024 * 1024) // 32MiB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(64 * 1024 * 1024)   // 32MiB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("breezy"),
		Description:                    &ruleDescription,
		RpsSoftLimitPerIp:              &rpsSoftLimitPerIp,
		RpsHardLimitPerIp:              &rpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      &MappingSecurityRuleDefaultResponseCodeOnMaxRequests,
		MaxConnectionsPerIp:            &maxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: &bandwidthBpsLimitPerConnection,
		BandwidthLimitOnlyAfterBytes:   &bandwidthLimitOnlyAfterBytes,
		ResponseCodeOnMaxConnections:   &MappingSecurityRuleDefaultResponseCodeOnMaxConnections,
	}
}

func MappingSecurityRulePresetPermissive() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Easy‑going yet alert. Basic protection without slowing down everyday traffic.",
	)
	rpsSoftLimitPerIp := uint(16)
	rpsHardLimitPerIp := uint(24)
	maxConnectionsPerIp := uint(12)
	bandwidthBpsLimitPerConnection := valueObject.Byte(16 * 1024 * 1024) // 16MiB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(32 * 1024 * 1024)   // 24MiB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("permissive"),
		Description:                    &ruleDescription,
		RpsSoftLimitPerIp:              &rpsSoftLimitPerIp,
		RpsHardLimitPerIp:              &rpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      &MappingSecurityRuleDefaultResponseCodeOnMaxRequests,
		MaxConnectionsPerIp:            &maxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: &bandwidthBpsLimitPerConnection,
		BandwidthLimitOnlyAfterBytes:   &bandwidthLimitOnlyAfterBytes,
		ResponseCodeOnMaxConnections:   &MappingSecurityRuleDefaultResponseCodeOnMaxConnections,
	}
}

func MappingSecurityRulePresetReasonable() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Sweet spot: prevent most abuses while minimizing false positives.",
	)
	rpsSoftLimitPerIp := uint(12)
	rpsHardLimitPerIp := uint(16)
	maxConnectionsPerIp := uint(8)
	bandwidthBpsLimitPerConnection := valueObject.Byte(8 * 1024 * 1024) // 8MiB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(16 * 1024 * 1024)  // 16MiB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("reasonable"),
		Description:                    &ruleDescription,
		RpsSoftLimitPerIp:              &rpsSoftLimitPerIp,
		RpsHardLimitPerIp:              &rpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      &MappingSecurityRuleDefaultResponseCodeOnMaxRequests,
		MaxConnectionsPerIp:            &maxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: &bandwidthBpsLimitPerConnection,
		BandwidthLimitOnlyAfterBytes:   &bandwidthLimitOnlyAfterBytes,
		ResponseCodeOnMaxConnections:   &MappingSecurityRuleDefaultResponseCodeOnMaxConnections,
	}
}

func MappingSecurityRulePresetVigilant() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Lower tolerances, higher standards - might be too restrictive.",
	)
	rpsSoftLimitPerIp := uint(8)
	rpsHardLimitPerIp := uint(12)
	maxConnectionsPerIp := uint(6)
	bandwidthBpsLimitPerConnection := valueObject.Byte(4 * 1024 * 1024) // 4MiB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(8 * 1024 * 1024)   // 12MiB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("vigilant"),
		Description:                    &ruleDescription,
		RpsSoftLimitPerIp:              &rpsSoftLimitPerIp,
		RpsHardLimitPerIp:              &rpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      &MappingSecurityRuleDefaultResponseCodeOnMaxRequests,
		MaxConnectionsPerIp:            &maxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: &bandwidthBpsLimitPerConnection,
		BandwidthLimitOnlyAfterBytes:   &bandwidthLimitOnlyAfterBytes,
		ResponseCodeOnMaxConnections:   &MappingSecurityRuleDefaultResponseCodeOnMaxConnections,
	}
}

func MappingSecurityRulePresetIronclad() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Draconian limits, permanent red alert. Use with caution.",
	)
	rpsSoftLimitPerIp := uint(6)
	rpsHardLimitPerIp := uint(8)
	maxConnectionsPerIp := uint(4)
	bandwidthBpsLimitPerConnection := valueObject.Byte(2 * 1024 * 1024) // 2MiB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(4 * 1024 * 1024)   // 8MiB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("ironclad"),
		Description:                    &ruleDescription,
		RpsSoftLimitPerIp:              &rpsSoftLimitPerIp,
		RpsHardLimitPerIp:              &rpsHardLimitPerIp,
		ResponseCodeOnMaxRequests:      &MappingSecurityRuleDefaultResponseCodeOnMaxRequests,
		MaxConnectionsPerIp:            &maxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: &bandwidthBpsLimitPerConnection,
		BandwidthLimitOnlyAfterBytes:   &bandwidthLimitOnlyAfterBytes,
		ResponseCodeOnMaxConnections:   &MappingSecurityRuleDefaultResponseCodeOnMaxConnections,
	}
}

func MappingSecurityRuleInitialPresets() []MappingSecurityRule {
	return []MappingSecurityRule{
		MappingSecurityRulePresetBreezy(), MappingSecurityRulePresetPermissive(),
		MappingSecurityRulePresetReasonable(), MappingSecurityRulePresetVigilant(),
		MappingSecurityRulePresetIronclad(),
	}
}
