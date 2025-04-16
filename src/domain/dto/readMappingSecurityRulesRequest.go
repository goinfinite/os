package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMappingSecurityRulesRequest struct {
	Pagination              Pagination                           `json:"pagination"`
	MappingSecurityRuleId   *valueObject.MappingSecurityRuleId   `json:"mappingSecurityRuleId"`
	MappingSecurityRuleName *valueObject.MappingSecurityRuleName `json:"mappingSecurityRuleName"`
	AllowedIp               *valueObject.IpAddress               `json:"allowedIp"`
	BlockedIp               *valueObject.IpAddress               `json:"blockedIp"`
	CreatedBeforeAt         *valueObject.UnixTime                `json:"createdBeforeAt"`
	CreatedAfterAt          *valueObject.UnixTime                `json:"createdAfterAt"`
}

type ReadMappingSecurityRulesResponse struct {
	Pagination           Pagination                   `json:"pagination"`
	MappingSecurityRules []entity.MappingSecurityRule `json:"mappingSecurityRules"`
}
