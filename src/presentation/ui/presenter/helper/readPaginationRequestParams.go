package uiPresenterHelper

import (
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	"github.com/labstack/echo/v4"
)

func ReadPaginationRequestParams(echoContext echo.Context) map[string]any {
	requestParams := map[string]any{}

	pageNumberParam := echoContext.QueryParam("pageNumber")
	pageNumber, pageNumberErr := tkVoUtil.InterfaceToUint32(pageNumberParam)
	if pageNumberErr == nil {
		requestParams["pageNumber"] = pageNumber
	}

	itemsPerPageParam := echoContext.QueryParam("itemsPerPage")
	itemsPerPage, itemsPerPageErr := tkVoUtil.InterfaceToUint16(itemsPerPageParam)
	if itemsPerPageErr == nil && itemsPerPage > 0 {
		requestParams["itemsPerPage"] = itemsPerPage
	}

	sortByParam := echoContext.QueryParam("sortBy")
	sortBy, sortByErr := tkValueObject.NewPaginationSortBy(sortByParam)
	if sortByErr == nil {
		requestParams["sortBy"] = sortBy.String()
	}

	sortDirectionParam := echoContext.QueryParam("sortDirection")
	sortDirection, sortDirectionErr := tkValueObject.NewPaginationSortDirection(
		sortDirectionParam,
	)
	if sortDirectionErr == nil {
		requestParams["sortDirection"] = sortDirection.String()
	}

	return requestParams
}
