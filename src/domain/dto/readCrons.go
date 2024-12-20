package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadCronsRequest struct {
	Pagination  Pagination               `json:"pagination"`
	CronId      *valueObject.CronId      `json:"id,omitempty"`
	CronComment *valueObject.CronComment `json:"comment,omitempty"`
}

type ReadCronsResponse struct {
	Pagination Pagination    `json:"pagination"`
	Crons      []entity.Cron `json:"crons"`
}
