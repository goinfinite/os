package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Copy(copyUnixFile dto.CopyUnixFile) error
	Compress(compressUnixFiles dto.CompressUnixFiles) (dto.CompressionProcessReport, error)
	Create(createUnixFile dto.CreateUnixFile) error
	Delete(unixFilePath valueObject.UnixFilePath) error
	Extract(extractUnixFiles dto.ExtractUnixFiles) error
	Move(
		unixSrcFilePath valueObject.UnixFilePath,
		unixDestinationPath valueObject.UnixFilePath,
		shouldOverwrite bool,
	) error
	UpdateContent(updateUnixFile dto.UpdateUnixFile) error
	UpdatePermissions(
		unixFilePath valueObject.UnixFilePath,
		unixFilePermissions valueObject.UnixFilePermissions,
	) error
	Upload(uploadUnixFiles dto.UploadUnixFiles) (dto.UploadProcessReport, error)
}
