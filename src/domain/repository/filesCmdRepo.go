package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Create(dto.AddUnixFile) error
	Move(dto.UpdateUnixFile) error
	Copy(dto.CopyUnixFile) error
	UpdateContent(dto.UpdateUnixFileContent) error
	UpdatePermissions(
		valueObject.UnixFilePath,
		valueObject.UnixFilePermissions,
	) error
	Delete([]valueObject.UnixFilePath)
	Compress(dto.CompressUnixFiles) (dto.CompressionProcessReport, error)
	Extract(valueObject.UnixFilePath, valueObject.UnixFilePath) error
	Upload(valueObject.UnixFilePath, valueObject.FileStreamHandler) error
}
