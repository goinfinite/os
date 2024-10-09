package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const paginationSortByRegex string = `^[\p{L}\d\.\-\ ]{1,256}$`

type PaginationSortBy string

func NewPaginationSortBy(value interface{}) (sortBy PaginationSortBy, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sortBy, errors.New("PaginationSortByMustBeString")
	}

	re := regexp.MustCompile(paginationSortByRegex)
	if !re.MatchString(stringValue) {
		return sortBy, errors.New("InvalidPaginationSortBy")
	}

	return PaginationSortBy(stringValue), nil
}

func (vo PaginationSortBy) String() string {
	return string(vo)
}
