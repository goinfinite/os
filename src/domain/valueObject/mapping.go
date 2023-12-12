package valueObject

type Mapping struct {
	Id                     MappingId           `json:"id"`
	Path                   UrlPath             `json:"path"`
	MatchPattern           MappingMatchPattern `json:"matchPattern"`
	TargetType             MappingTargetType   `json:"targetType"`
	TargetService          *ServiceName        `json:"targetService"`
	TargetUrl              *Url                `json:"targetUrl"`
	TargetHttpResponseCode *HttpResponseCode   `json:"targetHttpResponseCode"`
}

func NewMapping(
	id MappingId,
	path UrlPath,
	matchPattern MappingMatchPattern,
	targetType MappingTargetType,
	targetService *ServiceName,
	targetUrl *Url,
	targetHttpResponseCode *HttpResponseCode,
) Mapping {
	return Mapping{
		Id:                     id,
		Path:                   path,
		MatchPattern:           matchPattern,
		TargetType:             targetType,
		TargetService:          targetService,
		TargetUrl:              targetUrl,
		TargetHttpResponseCode: targetHttpResponseCode,
	}
}
