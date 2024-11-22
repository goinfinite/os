package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type MoveUnixFile struct {
	SourcePath      valueObject.UnixFilePath `json:"sourcePath"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
	ShouldOverwrite bool                     `json:"shouldOverwrite"`
}

func NewMoveUnixFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
	shouldOverwrite bool,
) MoveUnixFile {
	return MoveUnixFile{
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
		ShouldOverwrite: shouldOverwrite,
	}
}
