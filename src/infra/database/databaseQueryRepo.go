package databaseInfra

import (
	"errors"
	"math"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type DatabaseQueryRepo struct {
	dbType valueObject.DatabaseType
}

func NewDatabaseQueryRepo(
	dbType valueObject.DatabaseType,
) *DatabaseQueryRepo {
	return &DatabaseQueryRepo{dbType: dbType}
}

func (repo DatabaseQueryRepo) Read(
	requestDto dto.ReadDatabasesRequest,
) (responseDto dto.ReadDatabasesResponse, err error) {
	if requestDto.DatabaseType == nil {
		requestDto.DatabaseType = &repo.dbType
	}

	allDatabases := []entity.Database{}
	switch repo.dbType {
	case "mariadb":
		allDatabases, err = MysqlDatabaseQueryRepo{}.readAllDatabases()
	case "postgresql":
		allDatabases, err = PostgresDatabaseQueryRepo{}.readAllDatabases()
	default:
		return responseDto, errors.New("DatabaseTypeNotSupported")
	}
	if err != nil {
		return responseDto, err
	}

	filteredDatabases := []entity.Database{}
	for _, databaseEntity := range allDatabases {
		if requestDto.DatabaseName != nil && databaseEntity.Name != *requestDto.DatabaseName {
			continue
		}

		if requestDto.DatabaseType != nil && databaseEntity.Type != *requestDto.DatabaseType {
			continue
		}

		if requestDto.Username != nil {
			userFound := false
			for _, dbUserEntity := range databaseEntity.Users {
				if dbUserEntity.Username != *requestDto.Username {
					continue
				}

				userFound = true
				break
			}

			if !userFound {
				continue
			}
		}

		filteredDatabases = append(filteredDatabases, databaseEntity)
	}

	paginatedDatabases := []entity.Database{}
	startIndex := int(requestDto.Pagination.PageNumber) * int(requestDto.Pagination.ItemsPerPage)
	endIndex := startIndex + int(requestDto.Pagination.ItemsPerPage)

	if startIndex < len(filteredDatabases) {
		if endIndex > len(filteredDatabases) {
			endIndex = len(filteredDatabases)
		}
		paginatedDatabases = filteredDatabases[startIndex:endIndex]
	}

	itemsTotal := uint64(len(filteredDatabases))
	pagesTotal := uint32(math.Ceil(float64(itemsTotal) / float64(requestDto.Pagination.ItemsPerPage)))

	paginationDto := requestDto.Pagination
	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	responseDto.Pagination = paginationDto
	responseDto.Databases = paginatedDatabases

	return responseDto, nil
}

func (repo DatabaseQueryRepo) ReadFirst(
	requestDto dto.ReadDatabasesRequest,
) (database entity.Database, err error) {
	requestDto.Pagination = dto.PaginationSingleItem

	if requestDto.DatabaseType == nil {
		requestDto.DatabaseType = &repo.dbType
	}

	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return database, err
	}

	if len(responseDto.Databases) == 0 {
		return database, errors.New("DatabaseNotFound")
	}

	return responseDto.Databases[0], nil
}
