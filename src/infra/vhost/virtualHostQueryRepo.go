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

	dbQuery := repo.persistentDbSvc.Handler.Model(virtualHostModel).Where(&virtualHostModel)
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

	virtualHostModels := []dbModel.VirtualHost{}
	err = paginatedDbQuery.Find(&virtualHostModels).Error
	if err != nil {
		return responseDto, errors.New("FindVirtualHostsError: " + err.Error())
	}

	virtualHostEntities := []entity.VirtualHost{}
	for _, virtualHostModel := range virtualHostModels {
		virtualHostEntity, err := virtualHostModel.ToEntity()
		if err != nil {
			slog.Debug(
				"VirtualHostModelToEntityError",
				slog.String("hostname", virtualHostModel.Hostname),
				slog.String("error", err.Error()),
			)
			continue
		}
		virtualHostEntities = append(virtualHostEntities, virtualHostEntity)
	}

	if requestDto.WithMappings != nil && *requestDto.WithMappings {
		return dto.ReadVirtualHostsResponse{
			Pagination:   responsePagination,
			VirtualHosts: virtualHostEntities,
		}, nil
	}

	mappingQueryRepo := NewMappingQueryRepo(repo.persistentDbSvc)
	readMappingsResponse, err := mappingQueryRepo.Read(dto.ReadMappingsRequest{
		Hostname: requestDto.Hostname,
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
	for _, virtualHostEntity := range virtualHostEntities {
		if _, mappingExists := hostnameMappingEntitiesMap[virtualHostEntity.Hostname]; !mappingExists {
			continue
		}

		virtualHostWithMappings = append(virtualHostWithMappings, dto.VirtualHostWithMappings{
			VirtualHost: virtualHostEntity,
			Mappings:    hostnameMappingEntitiesMap[virtualHostEntity.Hostname],
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
		return mappingsFilePath, errors.New("VirtualHostNotFound")
	}

	if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
		if vhostEntity.ParentHostname == nil {
			return mappingsFilePath, errors.New("AliasMissingParentHostname")
		}

		vhostHostname = *vhostEntity.ParentHostname
	}

	isPrimary := true
	primaryVirtualHostEntity, err := repo.ReadFirst(dto.ReadVirtualHostsRequest{
		IsPrimary: &isPrimary,
	})
	if err != nil {
		return mappingsFilePath, errors.New("ReadPrimaryVirtualHostError: " + err.Error())
	}

	vhostFileNameStr := vhostHostname.String() + ".conf"
	if vhostEntity.Hostname == primaryVirtualHostEntity.Hostname {
		vhostFileNameStr = infraEnvs.PrimaryVirtualHostFileName
	}

	return valueObject.NewUnixFilePath(
		infraEnvs.MappingsConfDir + "/" + vhostFileNameStr,
	)
}
