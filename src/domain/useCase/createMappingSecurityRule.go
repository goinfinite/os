package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateMappingSecurityRule(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateMappingSecurityRule,
) (mappingSecurityRuleId valueObject.MappingSecurityRuleId, err error) {
	if createDto.ResponseCodeOnMaxRequests == nil {
		createDto.ResponseCodeOnMaxRequests = &entity.MappingSecurityRuleDefaultResponseCodeOnMaxRequests
	}
	if createDto.ResponseCodeOnMaxConnections == nil {
		createDto.ResponseCodeOnMaxConnections = &entity.MappingSecurityRuleDefaultResponseCodeOnMaxConnections
	}

	mappingSecurityRuleId, err = mappingCmdRepo.CreateSecurityRule(createDto)
	if err != nil {
		slog.Error("CreateMappingSecurityRuleError", slog.String("err", err.Error()))
		return mappingSecurityRuleId, errors.New("CreateMappingSecurityRuleInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateMappingSecurityRule(createDto, mappingSecurityRuleId)

	return mappingSecurityRuleId, nil
}
