package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadMappingSecurityRulesRequest struct {
	Pagination              tkDto.Pagination                     `json:"pagination"`
	MappingSecurityRuleId   *valueObject.MappingSecurityRuleId   `json:"mappingSecurityRuleId"`
	MappingSecurityRuleName *valueObject.MappingSecurityRuleName `json:"mappingSecurityRuleName"`
	AllowedIp               *tkValueObject.CidrBlock             `json:"allowedIp"`
	BlockedIp               *tkValueObject.CidrBlock             `json:"blockedIp"`
	CreatedBeforeAt         *tkValueObject.UnixTime              `json:"createdBeforeAt"`
	CreatedAfterAt          *tkValueObject.UnixTime              `json:"createdAfterAt"`
}

type ReadMappingSecurityRulesResponse struct {
	Pagination           tkDto.Pagination             `json:"pagination"`
	MappingSecurityRules []entity.MappingSecurityRule `json:"mappingSecurityRules"`
}
