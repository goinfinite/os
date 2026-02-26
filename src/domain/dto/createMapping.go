package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateMapping struct {
	Hostname                      tkValueObject.Fqdn                 `json:"hostname"`
	Path                          valueObject.MappingPath            `json:"path"`
	MatchPattern                  valueObject.MappingMatchPattern    `json:"matchPattern"`
	TargetType                    valueObject.MappingTargetType      `json:"targetType"`
	TargetValue                   *valueObject.MappingTargetValue    `json:"targetValue"`
	TargetHttpResponseCode        *tkValueObject.HttpStatusCode      `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool                              `json:"shouldUpgradeInsecureRequests"`
	MappingSecurityRuleId         *valueObject.MappingSecurityRuleId `json:"mappingSecurityRuleId"`
	OperatorAccountId             tkValueObject.AccountId            `json:"-"`
	OperatorIpAddress             tkValueObject.IpAddress            `json:"-"`
}

func NewCreateMapping(
	hostname tkValueObject.Fqdn,
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetValue *valueObject.MappingTargetValue,
	targetHttpResponseCode *tkValueObject.HttpStatusCode,
	shouldUpgradeInsecureRequests *bool,
	mappingSecurityRuleId *valueObject.MappingSecurityRuleId,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
