package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var MappingSecurityRulesDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadMappingSecurityRules(
	mappingQueryRepo repository.MappingQueryRepo,
	requestDto dto.ReadMappingSecurityRulesRequest,
) (responseDto dto.ReadMappingSecurityRulesResponse, err error) {
	responseDto, err = mappingQueryRepo.ReadSecurityRule(requestDto)
	if err != nil {
		slog.Error("ReadMappingSecurityRulesError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadMappingSecurityRulesInfraError")
	}

	return responseDto, nil
}
