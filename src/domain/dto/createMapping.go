package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateMapping struct {
	Hostname                      valueObject.Fqdn                   `json:"hostname"`
	Path                          valueObject.MappingPath            `json:"path"`
	MatchPattern                  valueObject.MappingMatchPattern    `json:"matchPattern"`
	TargetType                    valueObject.MappingTargetType      `json:"targetType"`
	TargetValue                   *valueObject.MappingTargetValue    `json:"targetValue"`
	TargetHttpResponseCode        *valueObject.HttpResponseCode      `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool                              `json:"shouldUpgradeInsecureRequests"`
	MappingSecurityRuleId         *valueObject.MappingSecurityRuleId `json:"mappingSecurityRuleId"`
	OperatorAccountId             valueObject.AccountId              `json:"-"`
	OperatorIpAddress             valueObject.IpAddress              `json:"-"`
}

func NewCreateMapping(
	hostname valueObject.Fqdn,
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetValue *valueObject.MappingTargetValue,
	targetHttpResponseCode *valueObject.HttpResponseCode,
	shouldUpgradeInsecureRequests *bool,
	mappingSecurityRuleId *valueObject.MappingSecurityRuleId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateMapping {
	return CreateMapping{
		Hostname:                      hostname,
		Path:                          path,
		MatchPattern:                  matchPattern,
		TargetType:                    targetType,
		TargetValue:                   targetValue,
		TargetHttpResponseCode:        targetHttpResponseCode,
		ShouldUpgradeInsecureRequests: shouldUpgradeInsecureRequests,
		MappingSecurityRuleId:         mappingSecurityRuleId,
		OperatorAccountId:             operatorAccountId,
		OperatorIpAddress:             operatorIpAddress,
	}
}
