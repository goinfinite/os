package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Add(dto.AddUnixFile) error
	UpdatePermissions(
		valueObject.UnixFilePath,
		valueObject.UnixFilePermissions,
		valueObject.UnixFileType,
	) error
}
