package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadMappingsRequest struct {
	Pagination                    Pagination                         `json:"pagination"`
	MappingId                     *valueObject.MappingId             `json:"mappingId"`
	Hostname                      *valueObject.Fqdn                  `json:"hostname"`
	MappingPath                   *valueObject.MappingPath           `json:"mappingPath"`
	MatchPattern                  *valueObject.MappingMatchPattern   `json:"matchPattern"`
	TargetType                    *valueObject.MappingTargetType     `json:"targetType"`
	TargetValue                   *valueObject.MappingTargetValue    `json:"targetValue"`
	TargetHttpResponseCode        *valueObject.HttpResponseCode      `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool                              `json:"shouldUpgradeInsecureRequests"`
	MappingSecurityRuleId         *valueObject.MappingSecurityRuleId `json:"mappingSecurityRuleId"`
	CreatedBeforeAt               *valueObject.UnixTime              `json:"createdBeforeAt"`
	CreatedAfterAt                *valueObject.UnixTime              `json:"createdAfterAt"`
}

type ReadMappingsResponse struct {
	Pagination Pagination       `json:"pagination"`
	Mappings   []entity.Mapping `json:"mappings"`
}
