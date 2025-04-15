package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type UpdateMappingSecurityRule struct {
	Id                             valueObject.MappingSecurityRuleId           `json:"id"`
	Name                           *valueObject.MappingSecurityRuleName        `json:"name"`
	Description                    *valueObject.MappingSecurityRuleDescription `json:"description"`
	AllowedIps                     *[]valueObject.IpAddress                    `json:"allowedIps"`
	BlockedIps                     *[]valueObject.IpAddress                    `json:"blockedIps"`
	SoftLimitRequestsPerIp         *uint                                       `json:"softLimitRequestsPerIp"`
	HardLimitRequestsPerIp         *uint                                       `json:"hardLimitRequestsPerIp"`
	ResponseCodeOnMaxRequests      *uint                                       `json:"responseCodeOnMaxRequests"`
	MaxConnectionsPerIp            *uint                                       `json:"maxConnectionsPerIp"`
	BandwidthBpsLimitPerConnection *valueObject.Byte                           `json:"bandwidthBpsLimitPerConnection"`
	BandwidthLimitOnlyAfterBytes   *valueObject.Byte                           `json:"bandwidthLimitOnlyAfterBytes"`
	ResponseCodeOnMaxConnections   *uint                                       `json:"responseCodeOnMaxConnections"`
	OperatorAccountId              valueObject.AccountId                       `json:"-"`
	OperatorIpAddress              valueObject.IpAddress                       `json:"-"`
}

func NewUpdateMappingSecurityRule(
	id valueObject.MappingSecurityRuleId,
	name *valueObject.MappingSecurityRuleName,
	description *valueObject.MappingSecurityRuleDescription,
	allowedIps, blockedIps *[]valueObject.IpAddress,
	softLimitRequestsPerIp, hardLimitRequestsPerIp, responseCodeOnMaxRequests, maxConnectionsPerIp *uint,
	bandwidthBpsLimitPerConnection, bandwidthLimitOnlyAfterBytes *valueObject.Byte,
	responseCodeOnMaxConnections *uint,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateMappingSecurityRule {
	return UpdateMappingSecurityRule{
		Id:                             id,
		Name:                           name,
		Description:                    description,
		AllowedIps:                     allowedIps,
		BlockedIps:                     blockedIps,
		SoftLimitRequestsPerIp:         softLimitRequestsPerIp,
		HardLimitRequestsPerIp:         hardLimitRequestsPerIp,
		ResponseCodeOnMaxRequests:      responseCodeOnMaxRequests,
		MaxConnectionsPerIp:            maxConnectionsPerIp,
		BandwidthBpsLimitPerConnection: bandwidthBpsLimitPerConnection,
		BandwidthLimitOnlyAfterBytes:   bandwidthLimitOnlyAfterBytes,
		ResponseCodeOnMaxConnections:   responseCodeOnMaxConnections,
		OperatorAccountId:              operatorAccountId,
		OperatorIpAddress:              operatorIpAddress,
	}
}
