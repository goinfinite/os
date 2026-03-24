package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateMappingSecurityRule struct {
	Name                           valueObject.MappingSecurityRuleName         `json:"name"`
	Description                    *valueObject.MappingSecurityRuleDescription `json:"description"`
	AllowedIps                     []tkValueObject.CidrBlock                   `json:"allowedIps"`
	BlockedIps                     []tkValueObject.CidrBlock                   `json:"blockedIps"`
	RpsSoftLimitPerIp              *uint                                       `json:"rpsSoftLimitPerIp"`
	RpsHardLimitPerIp              *uint                                       `json:"rpsHardLimitPerIp"`
	ResponseCodeOnMaxRequests      *uint                                       `json:"responseCodeOnMaxRequests"`
	MaxConnectionsPerIp            *uint                                       `json:"maxConnectionsPerIp"`
	BandwidthBpsLimitPerConnection *tkValueObject.Byte                         `json:"bandwidthBpsLimitPerConnection"`
	BandwidthLimitOnlyAfterBytes   *tkValueObject.Byte                         `json:"bandwidthLimitOnlyAfterBytes"`
	ResponseCodeOnMaxConnections   *uint                                       `json:"responseCodeOnMaxConnections"`
	OperatorAccountId              tkValueObject.AccountId                     `json:"-"`
	OperatorIpAddress              tkValueObject.IpAddress                     `json:"-"`
}

func NewCreateMappingSecurityRule(
	name valueObject.MappingSecurityRuleName,
	description *valueObject.MappingSecurityRuleDescription,
	allowedIps, blockedIps []tkValueObject.CidrBlock,
	rpsSoftLimitPerIp, rpsHardLimitPerIp, responseCodeOnMaxRequests, maxConnectionsPerIp *uint,
	bandwidthBpsLimitPerConnection, bandwidthLimitOnlyAfterBytes *tkValueObject.Byte,
	responseCodeOnMaxConnections *uint,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateMappingSecurityRule {
	return CreateMappingSecurityRule{
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
		OperatorAccountId:              operatorAccountId,
		OperatorIpAddress:              operatorIpAddress,
	}
}
