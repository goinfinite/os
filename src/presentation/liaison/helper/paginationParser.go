package liaisonHelper

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

func PaginationParser(
	untrustedInput map[string]any,
	defaultPagination dto.Pagination,
) (requestPagination dto.Pagination, err error) {
	requestPagination = defaultPagination

	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return requestPagination, errors.New("InvalidPageNumber")
		}
		requestPagination.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return requestPagination, errors.New("InvalidItemsPerPage")
		}
		requestPagination.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return requestPagination, err
		}
		requestPagination.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(untrustedInput["sortDirection"])
		if err != nil {
			return requestPagination, err
		}
		requestPagination.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return requestPagination, err
		}
		requestPagination.LastSeenId = &lastSeenId
	}

	return requestPagination, nil
}
