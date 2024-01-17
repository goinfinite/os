package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CopyUnixFile struct {
	SourcePath      valueObject.UnixFilePath `json:"sourcePath"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
}

func NewCopyUnixFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) CopyUnixFile {
	return CopyUnixFile{
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
	}
}
