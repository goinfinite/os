package dto

import "github.com/goinfinite/os/src/domain/valueObject"

var PaginationUnpaginated = Pagination{
	PageNumber:   0,
	ItemsPerPage: 1000,
}

type Pagination struct {
	PageNumber    uint32                               `json:"pageNumber"`
	ItemsPerPage  uint16                               `json:"itemsPerPage"`
	SortBy        *valueObject.PaginationSortBy        `json:"sortBy,omitempty"`
	SortDirection *valueObject.PaginationSortDirection `json:"sortDirection,omitempty"`
	LastSeenId    *valueObject.PaginationLastSeenId    `json:"lastSeenId,omitempty"`
	PagesTotal    *uint32                              `json:"pagesTotal,omitempty"`
	ItemsTotal    *uint64                              `json:"itemsTotal,omitempty"`
}
