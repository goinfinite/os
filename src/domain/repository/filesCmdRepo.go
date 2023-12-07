package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Add(dto.AddUnixFile) error
	Move(dto.MoveUnixFile) error
	UpdateContent(dto.UpdateUnixFileContent) error
	UpdatePermissions(
		valueObject.UnixFilePath,
		valueObject.UnixFilePermissions,
	) error
}
