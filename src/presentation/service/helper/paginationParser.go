package serviceHelper

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

func PaginationParser(
	userInput map[string]interface{},
	defaultPagination dto.Pagination,
) (requestPagination dto.Pagination, err error) {
	requestPagination = defaultPagination

	if userInput["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(userInput["pageNumber"])
		if err != nil {
			return requestPagination, errors.New("InvalidPageNumber")
		}
		requestPagination.PageNumber = pageNumber
	}

	if userInput["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(userInput["itemsPerPage"])
		if err != nil {
			return requestPagination, errors.New("InvalidItemsPerPage")
		}
		requestPagination.ItemsPerPage = itemsPerPage
	}

	if userInput["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(userInput["sortBy"])
		if err != nil {
			return requestPagination, err
		}
		requestPagination.SortBy = &sortBy
	}

	if userInput["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(userInput["sortDirection"])
		if err != nil {
			return requestPagination, err
		}
		requestPagination.SortDirection = &sortDirection
	}

	if userInput["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(userInput["lastSeenId"])
		if err != nil {
			return requestPagination, err
		}
		requestPagination.LastSeenId = &lastSeenId
	}

	return requestPagination, nil
}
