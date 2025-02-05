package dbHelper

import (
	"errors"
	"math"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

func PaginationQueryBuilder(
	dbQuery *gorm.DB,
	requestPagination dto.Pagination,
) (paginatedDbQuery *gorm.DB, responsePagination dto.Pagination, err error) {
	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return paginatedDbQuery, responsePagination, errors.New("CountItemsTotalError: " + err.Error())
	}

	//lint:ignore SA4006 paginatedDbQuery is used in the return statement
	paginatedDbQuery = dbQuery.Limit(int(requestPagination.ItemsPerPage))
	if requestPagination.LastSeenId == nil {
		offset := int(requestPagination.PageNumber) * int(requestPagination.ItemsPerPage)
		paginatedDbQuery = dbQuery.Offset(offset)
	} else {
		paginatedDbQuery = dbQuery.Where("id > ?", requestPagination.LastSeenId.String())
	}
	if requestPagination.SortBy != nil {
		orderStatement := requestPagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if requestPagination.SortDirection != nil {
			orderStatement += " " + requestPagination.SortDirection.String()
		}

		paginatedDbQuery = dbQuery.Order(orderStatement)
	}

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(requestPagination.ItemsPerPage)),
	)
	responsePagination = dto.Pagination{
		PageNumber:    requestPagination.PageNumber,
		ItemsPerPage:  requestPagination.ItemsPerPage,
		SortBy:        requestPagination.SortBy,
		SortDirection: requestPagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	return paginatedDbQuery, responsePagination, nil
}
