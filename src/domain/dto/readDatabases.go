package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadDatabasesRequest struct {
	Pagination   Pagination                  `json:"pagination"`
	DatabaseName *valueObject.DatabaseName   `json:"name,omitempty"`
	DatabaseType *valueObject.DatabaseType   `json:"type,omitempty"`
	Username     *valueObject.DatabaseUsername `json:"username,omitempty"`
}

type ReadDatabasesResponse struct {
	Pagination Pagination        `json:"pagination"`
	Databases  []entity.Database `json:"databases"`
}
