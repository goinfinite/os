package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateMapping struct {
	Hostname               valueObject.Fqdn                `json:"hostname"`
	Path                   valueObject.MappingPath         `json:"path"`
	MatchPattern           valueObject.MappingMatchPattern `json:"matchPattern"`
	TargetType             valueObject.MappingTargetType   `json:"targetType"`
	TargetServiceName      *valueObject.ServiceName        `json:"targetServiceName"`
	TargetUrl              *valueObject.Url                `json:"targetUrl"`
	TargetHttpResponseCode *valueObject.HttpResponseCode   `json:"targetHttpResponseCode"`
}

func NewCreateMapping(
	hostname valueObject.Fqdn,
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetServiceName *valueObject.ServiceName,
	targetUrl *valueObject.Url,
	targetHttpResponseCode *valueObject.HttpResponseCode,
) CreateMapping {
	return CreateMapping{
		Hostname:               hostname,
		Path:                   path,
		MatchPattern:           matchPattern,
		TargetType:             targetType,
		TargetServiceName:      targetServiceName,
		TargetUrl:              targetUrl,
		TargetHttpResponseCode: targetHttpResponseCode,
	}
}
