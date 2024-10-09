package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CopyUnixFile struct {
	SourcePath      valueObject.UnixFilePath `json:"sourcePath"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
	ShouldOverwrite bool                     `json:"shouldOverwrite"`
}

func NewCopyUnixFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
	shouldOverwrite bool,
) CopyUnixFile {
	return CopyUnixFile{
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
		ShouldOverwrite: shouldOverwrite,
	}
}
