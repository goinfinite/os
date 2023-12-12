package dto

import (
	"github.com/speedianet/os/src/domain/valueObject"
	apiValueObject "github.com/speedianet/os/src/presentation/api/valueObject"
)

type UploadUnixFiles struct {
	DestinationPath valueObject.UnixFilePath       `json:"destinationPath"`
	MultipartFiles  []apiValueObject.MultipartFile `json:"multipartFiles"`
}

func NewUploadUnixFiles(
	destinationPath valueObject.UnixFilePath,
	multipartFiles []apiValueObject.MultipartFile,
) UploadUnixFiles {
	return UploadUnixFiles{
		DestinationPath: destinationPath,
		MultipartFiles:  multipartFiles,
	}
}
