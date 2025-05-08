package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type MappingQueryRepo interface {
	Read(dto.ReadMappingsRequest) (dto.ReadMappingsResponse, error)
	ReadFirst(dto.ReadMappingsRequest) (entity.Mapping, error)
	ReadSecurityRule(
		dto.ReadMappingSecurityRulesRequest,
	) (dto.ReadMappingSecurityRulesResponse, error)
	ReadFirstSecurityRule(
		dto.ReadMappingSecurityRulesRequest,
	) (entity.MappingSecurityRule, error)
}
