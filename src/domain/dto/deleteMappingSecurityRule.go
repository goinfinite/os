package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteMappingSecurityRule struct {
	SecurityRuleId    valueObject.MappingSecurityRuleId `json:"securityRuleId"`
	OperatorAccountId valueObject.AccountId             `json:"-"`
	OperatorIpAddress valueObject.IpAddress             `json:"-"`
}

func NewDeleteMappingSecurityRule(
	securityRuleId valueObject.MappingSecurityRuleId,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteMappingSecurityRule {
	return DeleteMappingSecurityRule{
		SecurityRuleId:    securityRuleId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
