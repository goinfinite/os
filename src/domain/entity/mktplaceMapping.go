package entity

import "github.com/speedianet/os/src/domain/valueObject"

type MktplaceMapping struct {
	Path              valueObject.MappingPath         `json:"path"`
	MatchPattern      valueObject.MappingMatchPattern `json:"matchPattern"`
	TargetServiceName *valueObject.ServiceName        `json:"targetServiceName"`
}

func NewMktplaceMapping(
	path valueObject.MappingPath,
	matchPattern valueObject.MappingMatchPattern,
	targetServiceName *valueObject.ServiceName,
) MktplaceMapping {
	return MktplaceMapping{
		Path:              path,
		MatchPattern:      matchPattern,
		TargetServiceName: targetServiceName,
	}
}
