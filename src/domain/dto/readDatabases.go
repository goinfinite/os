package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadDatabasesRequest struct {
	Pagination   Pagination                    `json:"pagination"`
	DatabaseName *valueObject.DatabaseName     `json:"name"`
	DatabaseType *valueObject.DatabaseType     `json:"type"`
	Username     *valueObject.DatabaseUsername `json:"username"`
}

type ReadDatabasesResponse struct {
	Pagination Pagination        `json:"pagination"`
	Databases  []entity.Database `json:"databases"`
}
