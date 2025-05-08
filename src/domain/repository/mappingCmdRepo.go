package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MappingCmdRepo interface {
	Create(dto.CreateMapping) (valueObject.MappingId, error)
	Update(dto.UpdateMapping) error
	Delete(valueObject.MappingId) error
	CreateSecurityRule(
		dto.CreateMappingSecurityRule,
	) (valueObject.MappingSecurityRuleId, error)
	UpdateSecurityRule(dto.UpdateMappingSecurityRule) error
	DeleteSecurityRule(valueObject.MappingSecurityRuleId) error
}
