package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesQueryRepo interface {
	Get(unixFilePath valueObject.UnixFilePath) ([]entity.UnixFile, error)
	Download(unixFilePath valueObject.UnixFilePath) error
}
