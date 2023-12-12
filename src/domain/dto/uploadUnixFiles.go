package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type UploadUnixFiles struct {
	DestinationPath valueObject.UnixFilePath    `json:"destinationPath"`
	MultipartFiles  []valueObject.MultipartFile `json:"multipartFiles"`
}

func NewUploadUnixFiles(
	destinationPath valueObject.UnixFilePath,
	multipartFiles []valueObject.MultipartFile,
) UploadUnixFiles {
	return UploadUnixFiles{
		DestinationPath: destinationPath,
		MultipartFiles:  multipartFiles,
	}
}
