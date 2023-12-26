package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo interface {
	Create(addUnixFile dto.AddUnixFile) error
	Move(updateUnixFile dto.UpdateUnixFile) error
	Copy(copyUnixFile dto.CopyUnixFile) error
	UpdateContent(updateUnixFileContent dto.UpdateUnixFileContent) error
	UpdatePermissions(
		unixFilePath valueObject.UnixFilePath,
		unixFilePermissions valueObject.UnixFilePermissions,
	) error
	Compress(compressUnixFiles dto.CompressUnixFiles) dto.CompressionProcessReport
	Extract(extractUnixFiles dto.ExtractUnixFiles) error
	Delete(unixFilePathList []valueObject.UnixFilePath)
	Upload(uploadUnixFiles dto.UploadUnixFiles) dto.UploadProcessReport
}
