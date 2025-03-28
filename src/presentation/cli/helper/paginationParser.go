package cliHelper

func PaginationParser(
	requestBody map[string]interface{},
	paginationPageNumberUint32 uint32,
	paginationItemsPerPageUint16 uint16,
	paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string,
) map[string]interface{} {
	if paginationPageNumberUint32 != 0 {
		requestBody["pageNumber"] = paginationPageNumberUint32
	}
	if paginationItemsPerPageUint16 != 0 {
		requestBody["itemsPerPage"] = paginationItemsPerPageUint16
	}
	if paginationSortByStr != "" {
		requestBody["sortBy"] = paginationSortByStr
	}
	if paginationSortDirectionStr != "" {
		requestBody["sortDirection"] = paginationSortDirectionStr
	}
	if paginationLastSeenIdStr != "" {
		requestBody["lastSeenId"] = paginationLastSeenIdStr
	}

	return requestBody
}
