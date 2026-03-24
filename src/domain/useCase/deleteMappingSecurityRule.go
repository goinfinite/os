package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func DeleteMappingSecurityRule(
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteMappingSecurityRule,
) error {
	_, err := mappingQueryRepo.ReadFirstSecurityRule(dto.ReadMappingSecurityRulesRequest{
		Pagination:            tkDto.PaginationUnpaginated,
		MappingSecurityRuleId: &deleteDto.SecurityRuleId,
	})
	if err != nil {
		return errors.New("MappingSecurityRuleNotFound")
	}

	mappingsResponse, err := mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Pagination:            tkDto.PaginationUnpaginated,
		MappingSecurityRuleId: &deleteDto.SecurityRuleId,
	})
	if err == nil && len(mappingsResponse.Mappings) > 0 {
		return errors.New("MappingSecurityRuleInUse")
	}

	err = mappingCmdRepo.DeleteSecurityRule(deleteDto.SecurityRuleId)
	if err != nil {
		slog.Error("DeleteMappingSecurityRuleError", slog.String("err", err.Error()))
		return errors.New("DeleteMappingSecurityRuleInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteMappingSecurityRule(deleteDto)

	return nil
}
