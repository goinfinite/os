package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesQueryRepo interface {
	Exists(unixFilePath valueObject.UnixFilePath) bool
	Get(unixFilePath valueObject.UnixFilePath) ([]entity.UnixFile, error)
}
