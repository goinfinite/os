package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateMapping struct {
	Id                            valueObject.MappingId              `json:"id"`
	Path                          *valueObject.MappingPath           `json:"path"`
	MatchPattern                  *valueObject.MappingMatchPattern   `json:"matchPattern"`
	TargetType                    *valueObject.MappingTargetType     `json:"targetType"`
	TargetValue                   *valueObject.MappingTargetValue    `json:"targetValue"`
	TargetHttpResponseCode        *tkValueObject.HttpStatusCode      `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool                              `json:"shouldUpgradeInsecureRequests"`
	MappingSecurityRuleId         *valueObject.MappingSecurityRuleId `json:"mappingSecurityRuleId"`
	ClearableFields               []string                           `json:"-"`
	OperatorAccountId             tkValueObject.AccountId            `json:"-"`
	OperatorIpAddress             tkValueObject.IpAddress            `json:"-"`
}

func NewUpdateMapping(
	id valueObject.MappingId,
	path *valueObject.MappingPath,
	matchPattern *valueObject.MappingMatchPattern,
	targetType *valueObject.MappingTargetType,
	targetValue *valueObject.MappingTargetValue,
	targetHttpResponseCode *tkValueObject.HttpStatusCode,
	shouldUpgradeInsecureRequests *bool,
	mappingSecurityRuleId *valueObject.MappingSecurityRuleId,
	clearableFields []string,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
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
