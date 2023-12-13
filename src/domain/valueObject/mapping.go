package valueObject

type Mapping struct {
	Id                     MappingId           `json:"id"`
	Path                   MappingPath         `json:"path"`
	MatchPattern           MappingMatchPattern `json:"matchPattern"`
	TargetType             MappingTargetType   `json:"targetType"`
	TargetService          *ServiceName        `json:"targetService,omitempty"`
	TargetUrl              *Url                `json:"targetUrl,omitempty"`
	TargetHttpResponseCode *HttpResponseCode   `json:"targetHttpResponseCode,omitempty"`
}

func NewMapping(
	id MappingId,
	path MappingPath,
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
