package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateMapping struct {
	Hostname               valueObject.Fqdn                `json:"hostname"`
	Path                   valueObject.MappingPath         `json:"path"`
	MatchPattern           valueObject.MappingMatchPattern `json:"matchPattern"`
	TargetType             valueObject.MappingTargetType   `json:"targetType"`
	TargetValue            *valueObject.MappingTargetValue `json:"targetValue"`
	TargetHttpResponseCode *valueObject.HttpResponseCode   `json:"targetHttpResponseCode"`
}

func NewCreateMapping(
	hostname valueObject.Fqdn,
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetValue *valueObject.MappingTargetValue,
	targetHttpResponseCode *valueObject.HttpResponseCode,
) CreateMapping {
	return CreateMapping{
		Hostname:               hostname,
		Path:                   path,
		MatchPattern:           matchPattern,
		TargetType:             targetType,
		TargetValue:            targetValue,
		TargetHttpResponseCode: targetHttpResponseCode,
	}
}
