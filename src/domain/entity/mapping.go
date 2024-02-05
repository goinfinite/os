package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Mapping struct {
	Id                     valueObject.MappingId           `json:"id"`
	Hostname               valueObject.Fqdn                `json:"hostname"`
	Path                   valueObject.MappingPath         `json:"path"`
	MatchPattern           valueObject.MappingMatchPattern `json:"matchPattern"`
	TargetType             valueObject.MappingTargetType   `json:"targetType"`
	TargetServiceName      *valueObject.ServiceName        `json:"targetServiceName"`
	TargetUrl              *valueObject.Url                `json:"targetUrl"`
	TargetHttpResponseCode *valueObject.HttpResponseCode   `json:"targetHttpResponseCode"`
}

func NewMapping(
	id valueObject.MappingId,
	hostname valueObject.Fqdn,
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetType valueObject.MappingTargetType,
	targetServiceName *valueObject.ServiceName,
	targetUrl *valueObject.Url,
	targetHttpResponseCode *valueObject.HttpResponseCode,
) Mapping {
	return Mapping{
		Id:                     id,
		Hostname:               hostname,
		Path:                   path,
		MatchPattern:           matchPattern,
		TargetType:             targetType,
		TargetServiceName:      targetServiceName,
		TargetUrl:              targetUrl,
		TargetHttpResponseCode: targetHttpResponseCode,
	}
}
