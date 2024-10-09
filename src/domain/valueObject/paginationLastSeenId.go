package valueObject

import (
	"errors"
	"regexp"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const paginationLastSeenIdRegex string = `^[\w\-]{1,256}$`

type PaginationLastSeenId string

func NewPaginationLastSeenId(
	value interface{},
) (lastSeenId PaginationLastSeenId, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return lastSeenId, errors.New("PaginationLastSeenIdMustBeString")
	}

	re := regexp.MustCompile(paginationLastSeenIdRegex)
	if !re.MatchString(stringValue) {
		return lastSeenId, errors.New("InvalidPaginationLastSeenId")
	}

	return PaginationLastSeenId(stringValue), nil
}

func (vo PaginationLastSeenId) String() string {
	return string(vo)
}
