package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateMappingSecurityRule(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateMappingSecurityRule,
) (mappingSecurityRuleId valueObject.MappingSecurityRuleId, err error) {
	mappingSecurityRuleId, err = mappingCmdRepo.CreateSecurityRule(createDto)
	if err != nil {
		slog.Error("CreateMappingSecurityRuleError", slog.String("err", err.Error()))
		return mappingSecurityRuleId, errors.New("CreateMappingSecurityRuleInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateMappingSecurityRule(createDto, mappingSecurityRuleId)

	return mappingSecurityRuleId, nil
}
