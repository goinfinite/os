package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadMappingsRequest struct {
	Pagination                    tkDto.Pagination                   `json:"pagination"`
	MappingId                     *valueObject.MappingId             `json:"mappingId"`
	Hostname                      *tkValueObject.Fqdn                `json:"hostname"`
	MappingPath                   *valueObject.MappingPath           `json:"mappingPath"`
	MatchPattern                  *valueObject.MappingMatchPattern   `json:"matchPattern"`
	TargetType                    *valueObject.MappingTargetType     `json:"targetType"`
	TargetValue                   *valueObject.MappingTargetValue    `json:"targetValue"`
	TargetHttpResponseCode        *tkValueObject.HttpStatusCode      `json:"targetHttpResponseCode"`
	ShouldUpgradeInsecureRequests *bool                              `json:"shouldUpgradeInsecureRequests"`
	MappingSecurityRuleId         *valueObject.MappingSecurityRuleId `json:"mappingSecurityRuleId"`
	CreatedBeforeAt               *tkValueObject.UnixTime            `json:"createdBeforeAt"`
	CreatedAfterAt                *tkValueObject.UnixTime            `json:"createdAfterAt"`
}

type ReadMappingsResponse struct {
	Pagination tkDto.Pagination `json:"pagination"`
	Mappings   []entity.Mapping `json:"mappings"`
}
