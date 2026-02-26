package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

type ReadCronsRequest struct {
	Pagination  tkDto.Pagination         `json:"pagination"`
	CronId      *valueObject.CronId      `json:"id,omitempty"`
	CronComment *valueObject.CronComment `json:"comment,omitempty"`
}

type ReadCronsResponse struct {
	Pagination tkDto.Pagination `json:"pagination"`
	Crons      []entity.Cron    `json:"crons"`
}
