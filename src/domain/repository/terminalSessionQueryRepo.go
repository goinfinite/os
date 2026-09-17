package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
)

type TerminalSessionQueryRepo interface {
	Read(requestDto dto.ReadTerminalSessionsRequest) (dto.ReadTerminalSessionsResponse, error)
	ReadFirst(requestDto dto.ReadTerminalSessionsRequest) (entity.TerminalSession, error)
}
