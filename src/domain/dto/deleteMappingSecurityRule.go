package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteMappingSecurityRule struct {
	SecurityRuleId    valueObject.MappingSecurityRuleId `json:"securityRuleId"`
	OperatorAccountId tkValueObject.AccountId           `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress           `json:"-"`
}

func NewDeleteMappingSecurityRule(
	securityRuleId valueObject.MappingSecurityRuleId,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteMappingSecurityRule {
	return DeleteMappingSecurityRule{
		SecurityRuleId:    securityRuleId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
