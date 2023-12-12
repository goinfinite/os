package dto

import "github.com/speedianet/os/src/domain/valueObject"

type ExtractUnixFiles struct {
	Path            valueObject.UnixFilePath `json:"path"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
}

func NewExtractUnixFiles(
	Path valueObject.UnixFilePath,
	DestinationPath valueObject.UnixFilePath,
) ExtractUnixFiles {
	return ExtractUnixFiles{
		Path:            Path,
		DestinationPath: DestinationPath,
	}
}
