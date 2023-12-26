package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesQueryRepo interface {
	IsDir(unixFilePath valueObject.UnixFilePath) (bool, error)
	Get(unixFilePath valueObject.UnixFilePath) ([]entity.UnixFile, error)
	GetOnly(unixFilePath valueObject.UnixFilePath) (entity.UnixFile, error)
}
