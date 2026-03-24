package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

type ReadDatabasesRequest struct {
	Pagination   tkDto.Pagination              `json:"pagination"`
	DatabaseName *valueObject.DatabaseName     `json:"name"`
	DatabaseType *valueObject.DatabaseType     `json:"type"`
	Username     *valueObject.DatabaseUsername `json:"username"`
}

type ReadDatabasesResponse struct {
	Pagination tkDto.Pagination  `json:"pagination"`
	Databases  []entity.Database `json:"databases"`
}
