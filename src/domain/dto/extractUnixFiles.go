package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type ExtractUnixFiles struct {
	SourcePath      valueObject.UnixFilePath `json:"sourcePath"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
}

func NewExtractUnixFiles(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) ExtractUnixFiles {
	return ExtractUnixFiles{
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
	}
}
