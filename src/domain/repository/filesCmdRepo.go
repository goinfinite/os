package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Add(dto.AddUnixFile) error
	Move(valueObject.UnixFilePath, valueObject.UnixFilePath) error
	UpdateContent(dto.UpdateUnixFileContent) error
	UpdatePermissions(
		valueObject.UnixFilePath,
		valueObject.UnixFilePermissions,
	) error
}
