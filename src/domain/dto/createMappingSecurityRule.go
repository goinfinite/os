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
	BandwidthBpsLimitPerConnection *valueObject.Byte                           `json:"bandwidthBpsLimitPerConnection"`
	BandwidthLimitOnlyAfterBytes   *valueObject.Byte                           `json:"bandwidthLimitOnlyAfterBytes"`
	ResponseCodeOnMaxConnections   *uint                                       `json:"responseCodeOnMaxConnections"`
	OperatorAccountId              valueObject.AccountId                       `json:"-"`
	OperatorIpAddress              valueObject.IpAddress                       `json:"-"`
}

func NewCreateMappingSecurityRule(
	name valueObject.MappingSecurityRuleName,
	description *valueObject.MappingSecurityRuleDescription,
	allowedIps, blockedIps []tkValueObject.CidrBlock,
	rpsSoftLimitPerIp, rpsHardLimitPerIp, responseCodeOnMaxRequests, maxConnectionsPerIp *uint,
	bandwidthBpsLimitPerConnection, bandwidthLimitOnlyAfterBytes *valueObject.Byte,
	responseCodeOnMaxConnections *uint,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
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
