package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type UploadUnixFiles struct {
	DestinationPath    valueObject.UnixFilePath        `json:"destinationPath"`
	FileStreamHandlers []valueObject.FileStreamHandler `json:"fileStreamHandlers"`
}

func NewUploadUnixFiles(
	destinationPath valueObject.UnixFilePath,
	fileStreamHandlers []valueObject.FileStreamHandler,
) UploadUnixFiles {
	return UploadUnixFiles{
		DestinationPath:    destinationPath,
		FileStreamHandlers: fileStreamHandlers,
	}
}
