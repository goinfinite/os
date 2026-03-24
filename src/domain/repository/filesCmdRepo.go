package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Copy(dto.CopyUnixFile) error
	Compress(dto.CompressUnixFiles) (dto.CompressionProcessReport, error)
	Create(dto.CreateUnixFile) error
	Delete(tkValueObject.UnixAbsoluteFilePath) error
	Extract(dto.ExtractUnixFiles) error
	Move(dto.MoveUnixFile) error
	UpdateContent(dto.UpdateUnixFileContent) error
	UpdateOwnership(dto.UpdateUnixFileOwnership) error
	UpdatePermissions(dto.UpdateUnixFilePermissions) error
	Upload(dto.UploadUnixFiles) (dto.UploadProcessReport, error)
}
