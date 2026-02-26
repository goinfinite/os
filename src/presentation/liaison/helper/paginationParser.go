package liaisonHelper

import (
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

func PaginationParser(
	untrustedInput map[string]any,
	defaultPagination tkDto.Pagination,
) (requestPagination tkDto.Pagination, err error) {
	return tkPresentation.PaginationParser(defaultPagination, untrustedInput)
}
