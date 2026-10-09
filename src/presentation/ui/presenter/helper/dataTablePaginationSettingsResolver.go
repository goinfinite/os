package uiPresenterHelper

import (
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	uiStructural "github.com/goinfinite/ui/src/structural"
)

type DataTablePaginationSettings struct {
	PageNumber    uint
	ItemsPerPage  uint
	ItemsTotal    uint
	PagesTotal    uint
	SortKey       string
	SortDirection uiStructural.DataTableSortDirection
}

func DataTablePaginationSettingsResolver(
	pagination tkDto.Pagination,
) DataTablePaginationSettings {
	settings := DataTablePaginationSettings{
		PageNumber:    uint(pagination.PageNumber),
		ItemsPerPage:  uint(pagination.ItemsPerPage),
		SortDirection: uiStructural.DataTableSortDirectionAsc,
	}

	if pagination.ItemsTotal != nil {
		settings.ItemsTotal = uint(*pagination.ItemsTotal)
	}
	if pagination.PagesTotal != nil {
		settings.PagesTotal = uint(*pagination.PagesTotal)
	}
	if pagination.SortBy != nil {
		settings.SortKey = pagination.SortBy.String()
	}
	if pagination.SortDirection != nil {
		settings.SortDirection = uiStructural.DataTableSortDirection(
			pagination.SortDirection.String(),
		)
	}

	return settings
}
