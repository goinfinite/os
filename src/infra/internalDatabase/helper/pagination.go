package dbHelper

import (
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
	"gorm.io/gorm"
)

func PaginationQueryBuilder(
	dbQuery *gorm.DB,
	requestPagination tkDto.Pagination,
) (paginatedDbQuery *gorm.DB, responsePagination tkDto.Pagination, err error) {
	return tkInfraDb.PaginationQueryBuilder(dbQuery, requestPagination)
}
