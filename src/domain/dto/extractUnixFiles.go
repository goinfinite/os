package dto

import "github.com/speedianet/os/src/domain/valueObject"

type ExtractUnixFiles struct {
	Path            valueObject.UnixFilePath `json:"path"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
}

func NewExtractUnixFiles(
	path valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) ExtractUnixFiles {
	return ExtractUnixFiles{
		Path:            path,
		DestinationPath: destinationPath,
	}
}
