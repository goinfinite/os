package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type FilesQueryRepo interface {
	Read(unixFilePath valueObject.UnixFilePath) ([]entity.UnixFile, error)
	ReadUnique(unixFilePath valueObject.UnixFilePath) (entity.UnixFile, error)
}
