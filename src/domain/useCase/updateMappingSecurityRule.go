package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func UpdateMappingSecurityRule(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	updateDto dto.UpdateMappingSecurityRule,
) error {
	_, err := mappingQueryRepo.ReadFirstSecurityRule(
		dto.ReadMappingSecurityRulesRequest{MappingSecurityRuleId: &updateDto.Id},
	)
	if err != nil {
		return errors.New("MappingSecurityRuleNotFound")
	}

	err = mappingCmdRepo.UpdateSecurityRule(updateDto)
	if err != nil {
		slog.Error("UpdateMappingSecurityRuleError", slog.String("err", err.Error()))
		return errors.New("UpdateMappingSecurityRuleInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		UpdateMappingSecurityRule(updateDto)

	return nil
}
