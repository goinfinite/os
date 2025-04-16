package entity

import "github.com/goinfinite/os/src/domain/valueObject"

var (
	MappingSecurityRuleDefaultResponseCodeOnMaxRequests    uint = 429
	MappingSecurityRuleDefaultResponseCodeOnMaxConnections uint = 420
)

type MappingSecurityRule struct {
	Id                             valueObject.MappingSecurityRuleId           `json:"id"`
	Name                           valueObject.MappingSecurityRuleName         `json:"name"`
	Description                    *valueObject.MappingSecurityRuleDescription `json:"description"`
	AllowedIps                     []valueObject.IpAddress                     `json:"allowedIps"`
	BlockedIps                     []valueObject.IpAddress                     `json:"blockedIps"`
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
	allowedIps, blockedIps []valueObject.IpAddress,
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

func MappingSecurityRulePresetRelaxed() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Basic restrictions, likely to go unnoticed by most.",
	)
	rpsSoftLimitPerIp := uint(32)
	rpsHardLimitPerIp := uint(64)
	maxConnectionsPerIp := uint(32)
	bandwidthBpsLimitPerConnection := valueObject.Byte(32 * 1024 * 1024) // 32MB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(64 * 1024 * 1024)   // 64MB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("relaxed"),
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

func MappingSecurityRulePresetLow() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Low restrictions, suitable for most use cases.",
	)
	rpsSoftLimitPerIp := uint(24)
	rpsHardLimitPerIp := uint(48)
	maxConnectionsPerIp := uint(24)
	bandwidthBpsLimitPerConnection := valueObject.Byte(24 * 1024 * 1024) // 24MB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(48 * 1024 * 1024)   // 48MB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("low"),
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

func MappingSecurityRulePresetMedium() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Medium restrictions, may limit some traffic patterns.",
	)
	rpsSoftLimitPerIp := uint(16)
	rpsHardLimitPerIp := uint(32)
	maxConnectionsPerIp := uint(16)
	bandwidthBpsLimitPerConnection := valueObject.Byte(16 * 1024 * 1024) // 16MB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(32 * 1024 * 1024)   // 32MB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("medium"),
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

func MappingSecurityRulePresetHigh() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"High restrictions, useful for resource-intensive endpoints.",
	)
	rpsSoftLimitPerIp := uint(12)
	rpsHardLimitPerIp := uint(16)
	maxConnectionsPerIp := uint(12)
	bandwidthBpsLimitPerConnection := valueObject.Byte(12 * 1024 * 1024) // 12MB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(24 * 1024 * 1024)   // 24MB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("high"),
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

func MappingSecurityRulePresetStrict() MappingSecurityRule {
	ruleDescription, _ := valueObject.NewMappingSecurityRuleDescription(
		"Strict restrictions, suitable for highly sensitive endpoints.",
	)
	rpsSoftLimitPerIp := uint(8)
	rpsHardLimitPerIp := uint(12)
	maxConnectionsPerIp := uint(8)
	bandwidthBpsLimitPerConnection := valueObject.Byte(8 * 1024 * 1024) // 8MB
	bandwidthLimitOnlyAfterBytes := valueObject.Byte(16 * 1024 * 1024)  // 16MB

	return MappingSecurityRule{
		Name:                           valueObject.MappingSecurityRuleName("strict"),
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
		MappingSecurityRulePresetRelaxed(), MappingSecurityRulePresetLow(),
		MappingSecurityRulePresetMedium(), MappingSecurityRulePresetHigh(),
		MappingSecurityRulePresetStrict(),
	}
}
