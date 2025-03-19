package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type FilesQueryRepo interface {
	Read(dto.ReadFilesRequest) (dto.ReadFilesResponse, error)
	ReadFirst(valueObject.UnixFilePath) (entity.UnixFile, error)
}
