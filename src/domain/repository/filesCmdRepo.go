package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Copy(copyUnixFile dto.CopyUnixFile) error
	Compress(compressUnixFiles dto.CompressUnixFiles) dto.CompressionProcessReport
	Create(createUnixFile dto.CreateUnixFile) error
	Delete(unixFilePathList []valueObject.UnixFilePath)
	Extract(extractUnixFiles dto.ExtractUnixFiles) error
	Move(updateUnixFile dto.UpdateUnixFile) error
	UpdateContent(updateUnixFileContent dto.UpdateUnixFileContent) error
	UpdatePermissions(
		unixFilePath valueObject.UnixFilePath,
		unixFilePermissions valueObject.UnixFilePermissions,
	) error
	Upload(uploadUnixFiles dto.UploadUnixFiles) dto.UploadProcessReport
}
