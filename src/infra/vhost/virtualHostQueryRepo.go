package vhostInfra

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbHelper "github.com/goinfinite/os/src/infra/internalDatabase/helper"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
)

type VirtualHostQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewVirtualHostQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostQueryRepo {
	return &VirtualHostQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *VirtualHostQueryRepo) Read(requestDto dto.ReadVirtualHostsRequest) (
	responseDto dto.ReadVirtualHostsResponse, err error,
) {
	virtualHostModel := dbModel.VirtualHost{}
	if requestDto.Hostname != nil {
		virtualHostModel.Hostname = requestDto.Hostname.String()
	}
	if requestDto.VirtualHostType != nil {
		virtualHostModel.Type = requestDto.VirtualHostType.String()
	}
	if requestDto.RootDirectory != nil {
		virtualHostModel.RootDirectory = requestDto.RootDirectory.String()
	}
	if requestDto.ParentHostname != nil {
		parentHostnameStr := requestDto.ParentHostname.String()
		virtualHostModel.ParentHostname = &parentHostnameStr
	}
	if requestDto.IsPrimary != nil && *requestDto.IsPrimary {
		primaryHostname, err := infraHelper.ReadPrimaryVirtualHostHostname()
		if err != nil {
			return responseDto, errors.New("ReadPrimaryVirtualHostHostnameError: " + err.Error())
		}
		virtualHostModel.Hostname = primaryHostname.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.Model(virtualHostModel).
		Preload("Aliases").Where(&virtualHostModel)
	if requestDto.IsWildcard != nil {
		dbQuery = dbQuery.Where("is_wildcard = ?", *requestDto.IsWildcard)
	}
	if len(requestDto.AliasesHostnames) > 0 {
		dbQuery = dbQuery.Where("Aliases.Hostname = ?", requestDto.AliasesHostnames)
	}
	if requestDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", requestDto.CreatedBeforeAt.ReadAsGoTime())
	}
	if requestDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", requestDto.CreatedAfterAt.ReadAsGoTime())
	}

	if requestDto.Pagination.SortBy != nil {
		sortByStr := requestDto.Pagination.SortBy.String()
		switch sortByStr {
		case "virtualHostType":
			sortByStr = "Type"
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

	vhostModels := []dbModel.VirtualHost{}
	err = paginatedDbQuery.Find(&vhostModels).Error
	if err != nil {
		return responseDto, errors.New("FindVirtualHostsError: " + err.Error())
	}

	vhostEntities := []entity.VirtualHost{}
	for _, virtualHostModel := range vhostModels {
		vhostEntity, err := virtualHostModel.ToEntity()
		if err != nil {
			slog.Debug(
				"VirtualHostModelToEntityError",
				slog.String("hostname", virtualHostModel.Hostname),
				slog.String("error", err.Error()),
			)
			continue
		}
		vhostEntities = append(vhostEntities, vhostEntity)
	}

	responseDto = dto.ReadVirtualHostsResponse{
		Pagination:   responsePagination,
		VirtualHosts: vhostEntities,
	}

	if requestDto.WithMappings == nil {
		return responseDto, nil
	}
	if !*requestDto.WithMappings {
		return responseDto, nil
	}

	mappingQueryRepo := NewMappingQueryRepo(repo.persistentDbSvc)
	readMappingsResponse, err := mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Pagination: dto.PaginationUnpaginated,
		Hostname:   requestDto.Hostname,
	})
	if err != nil {
		return responseDto, errors.New("ReadMappingsError: " + err.Error())
	}

	hostnameMappingEntitiesMap := map[valueObject.Fqdn][]entity.Mapping{}
	for _, mappingEntity := range readMappingsResponse.Mappings {
		hostnameMappingEntitiesMap[mappingEntity.Hostname] = append(
			hostnameMappingEntitiesMap[mappingEntity.Hostname], mappingEntity,
		)
	}

	virtualHostWithMappings := []dto.VirtualHostWithMappings{}
	for _, virtualHostEntity := range responseDto.VirtualHosts {
		mappingEntities := []entity.Mapping{}
		if _, mappingExists := hostnameMappingEntitiesMap[virtualHostEntity.Hostname]; mappingExists {
			mappingEntities = hostnameMappingEntitiesMap[virtualHostEntity.Hostname]
		}

		virtualHostWithMappings = append(virtualHostWithMappings, dto.VirtualHostWithMappings{
			VirtualHost: virtualHostEntity,
			Mappings:    mappingEntities,
		})
	}

	return dto.ReadVirtualHostsResponse{
		Pagination:              responsePagination,
		VirtualHostWithMappings: virtualHostWithMappings,
	}, nil
}

func (repo *VirtualHostQueryRepo) ReadFirst(
	requestDto dto.ReadVirtualHostsRequest,
) (vhostEntity entity.VirtualHost, err error) {
	requestDto.WithMappings = nil
	requestDto.Pagination = dto.PaginationSingleItem
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return vhostEntity, err
	}

	if len(responseDto.VirtualHosts) == 0 {
		return vhostEntity, errors.New("VirtualHostNotFound")
	}

	return responseDto.VirtualHosts[0], nil
}

func (repo *VirtualHostQueryRepo) ReadFirstWithMappings(
	requestDto dto.ReadVirtualHostsRequest,
) (vhostWithMappingsDto dto.VirtualHostWithMappings, err error) {
	withMappings := true
	requestDto.WithMappings = &withMappings
	requestDto.Pagination = dto.PaginationSingleItem

	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return vhostWithMappingsDto, err
	}

	if len(responseDto.VirtualHostWithMappings) == 0 {
		return vhostWithMappingsDto, errors.New("VirtualHostNotFound")
	}

	return responseDto.VirtualHostWithMappings[0], nil
}

func (repo *VirtualHostQueryRepo) ReadVirtualHostMappingsFilePath(
	vhostHostname valueObject.Fqdn,
) (mappingsFilePath valueObject.UnixFilePath, err error) {
	vhostEntity, err := repo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &vhostHostname,
	})
	if err != nil {
		return mappingsFilePath, errors.New("ReadVirtualHostEntityError: " + err.Error())
	}

	if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
		if vhostEntity.ParentHostname == nil {
			return mappingsFilePath, errors.New("AliasMissingParentHostname")
		}

		vhostEntity, err = repo.ReadFirst(dto.ReadVirtualHostsRequest{
			Hostname: vhostEntity.ParentHostname,
		})
		if err != nil {
			return mappingsFilePath, errors.New(
				"ReadParentVirtualHostEntityError: " + err.Error(),
			)
		}
		vhostHostname = vhostEntity.Hostname
	}

	vhostFileNameStr := vhostHostname.String() + ".conf"
	if vhostEntity.IsPrimary {
		vhostFileNameStr = "primary.conf"
	}

	return valueObject.NewUnixFilePath(
		infraEnvs.MappingsConfDir + "/" + vhostFileNameStr,
	)
}
