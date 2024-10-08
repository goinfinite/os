package repository

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type FilesQueryRepo interface {
	Get(unixFilePath valueObject.UnixFilePath) ([]entity.UnixFile, error)
	GetOne(unixFilePath valueObject.UnixFilePath) (entity.UnixFile, error)
}
