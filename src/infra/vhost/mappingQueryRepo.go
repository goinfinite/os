package vhostInfra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbHelper "github.com/goinfinite/os/src/infra/internalDatabase/helper"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
)

type MappingQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMappingQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MappingQueryRepo {
	return &MappingQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *MappingQueryRepo) Read(requestDto dto.ReadMappingsRequest) (
	responseDto dto.ReadMappingsResponse, err error,
) {
	mappingModel := dbModel.Mapping{}
	if requestDto.MappingId != nil {
		mappingModel.ID = requestDto.MappingId.Uint64()
	}
	if requestDto.Hostname != nil {
		mappingModel.Hostname = requestDto.Hostname.String()
	}
	if requestDto.MappingPath != nil {
		mappingModel.Path = requestDto.MappingPath.String()
	}
	if requestDto.MatchPattern != nil {
		mappingModel.MatchPattern = requestDto.MatchPattern.String()
	}
	if requestDto.TargetType != nil {
		mappingModel.TargetType = requestDto.TargetType.String()
	}
	if requestDto.TargetValue != nil {
		targetValueStr := requestDto.TargetValue.String()
		mappingModel.TargetValue = &targetValueStr
	}
	if requestDto.TargetHttpResponseCode != nil {
		targetHttpResponseCodeStr := requestDto.TargetHttpResponseCode.String()
		mappingModel.TargetHttpResponseCode = &targetHttpResponseCodeStr
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(mappingModel).Where(&mappingModel)
	if requestDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", requestDto.CreatedBeforeAt.ReadAsGoTime())
	}
	if requestDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", requestDto.CreatedAfterAt.ReadAsGoTime())
	}

	if requestDto.Pagination.SortBy != nil {
		sortByStr := requestDto.Pagination.SortBy.String()
		switch sortByStr {
		case "mappingId", "id":
			sortByStr = "ID"
		case "mappingPath":
			sortByStr = "Path"
		}

		sortBy, err := valueObject.NewPaginationSortBy(sortByStr)
		if err == nil {
			requestDto.Pagination.SortBy = &sortBy
		}
	}
	paginatedDbQuery, responsePagination, err := dbHelper.PaginationQueryBuilder(
		dbQuery, requestDto.Pagination,
	)
	if err != nil {
		return responseDto, errors.New("PaginationQueryBuilderError: " + err.Error())
	}

	mappingModels := []dbModel.Mapping{}
	err = paginatedDbQuery.Find(&mappingModels).Error
	if err != nil {
		return responseDto, errors.New("FindMappingsError: " + err.Error())
	}

	mappingEntities := []entity.Mapping{}
	for _, mappingModel := range mappingModels {
		mappingEntity, err := mappingModel.ToEntity()
		if err != nil {
			slog.Debug(
				"MappingModelToEntityError",
				slog.Uint64("mappingId", mappingModel.ID),
				slog.String("hostname", mappingModel.Hostname),
				slog.String("error", err.Error()),
			)
			continue
		}
		mappingEntities = append(mappingEntities, mappingEntity)
	}

	return dto.ReadMappingsResponse{
		Pagination: responsePagination,
		Mappings:   mappingEntities,
	}, nil
}

func (repo *MappingQueryRepo) ReadFirst(
	requestDto dto.ReadMappingsRequest,
) (mappingEntity entity.Mapping, err error) {
	requestDto.Pagination = dto.PaginationSingleItem
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return mappingEntity, err
	}

	if len(responseDto.Mappings) == 0 {
		return mappingEntity, errors.New("MappingNotFound")
	}

	return responseDto.Mappings[0], nil
}
