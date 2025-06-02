package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateMapping struct {
	Id                            valueObject.MappingId              `json:"id"`
	Path                          *valueObject.MappingPath           `json:"path"`
	MatchPattern                  *valueObject.MappingMatchPattern   `json:"matchPattern"`
	TargetType                    *valueObject.MappingTargetType     `json:"targetType"`
	TargetValue                   *valueObject.MappingTargetValue    `json:"targetValue"`
	TargetHttpResponseCode        *valueObject.HttpResponseCode      `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool                              `json:"shouldUpgradeInsecureRequests"`
	MappingSecurityRuleId         *valueObject.MappingSecurityRuleId `json:"mappingSecurityRuleId"`
	ClearableFields               []string                           `json:"-"`
	OperatorAccountId             valueObject.AccountId              `json:"-"`
	OperatorIpAddress             valueObject.IpAddress              `json:"-"`
}

func NewUpdateMapping(
	id valueObject.MappingId,
	path *valueObject.MappingPath,
	matchPattern *valueObject.MappingMatchPattern,
	targetType *valueObject.MappingTargetType,
	targetValue *valueObject.MappingTargetValue,
	targetHttpResponseCode *valueObject.HttpResponseCode,
	shouldUpgradeInsecureRequests *bool,
	mappingSecurityRuleId *valueObject.MappingSecurityRuleId,
	clearableFields []string,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateMapping {
	return UpdateMapping{
		Id:                            id,
		Path:                          path,
		MatchPattern:                  matchPattern,
		TargetType:                    targetType,
		TargetValue:                   targetValue,
		TargetHttpResponseCode:        targetHttpResponseCode,
		ShouldUpgradeInsecureRequests: shouldUpgradeInsecureRequests,
		MappingSecurityRuleId:         mappingSecurityRuleId,
		ClearableFields:               clearableFields,
		OperatorAccountId:             operatorAccountId,
		OperatorIpAddress:             operatorIpAddress,
	}
}
