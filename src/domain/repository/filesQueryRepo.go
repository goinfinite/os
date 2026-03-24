package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type FilesQueryRepo interface {
	Read(dto.ReadFilesRequest) (dto.ReadFilesResponse, error)
	ReadFirst(tkValueObject.UnixAbsoluteFilePath) (entity.UnixFile, error)
}
