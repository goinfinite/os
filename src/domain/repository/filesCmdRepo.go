package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Copy(dto.CopyUnixFile) error
	Compress(dto.CompressUnixFiles) (dto.CompressionProcessReport, error)
	Create(dto.CreateUnixFile) error
	Delete(valueObject.UnixFilePath) error
	Extract(dto.ExtractUnixFiles) error
	Move(dto.MoveUnixFile) error
	UpdateContent(dto.UpdateUnixFileContent) error
	UpdatePermissions(dto.UpdateUnixFilePermissions) error
	Upload(dto.UploadUnixFiles) (dto.UploadProcessReport, error)
}
