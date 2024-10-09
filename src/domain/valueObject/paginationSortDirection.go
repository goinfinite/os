package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type PaginationSortDirection string

var ValidPaginationSortDirections = []string{"asc", "desc"}

func NewPaginationSortDirection(
	value interface{},
) (sortDirection PaginationSortDirection, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sortDirection, errors.New("PaginationSortDirectionMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidPaginationSortDirections, stringValue) {
		return sortDirection, errors.New("InvalidPaginationSortDirection")
	}

	return PaginationSortDirection(stringValue), nil
}

func (vo PaginationSortDirection) String() string {
	return string(vo)
}
