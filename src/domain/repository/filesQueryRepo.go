package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type FilesQueryRepo interface {
	Read(valueObject.UnixFilePath) ([]entity.UnixFile, error)
	ReadFirst(valueObject.UnixFilePath) (entity.UnixFile, error)
}
