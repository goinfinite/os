package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type MoveUnixFile struct {
	SourcePath      tkValueObject.UnixAbsoluteFilePath `json:"sourcePath"`
	DestinationPath tkValueObject.UnixAbsoluteFilePath `json:"destinationPath"`
	ShouldOverwrite bool                               `json:"shouldOverwrite"`
}

func NewMoveUnixFile(
	sourcePath tkValueObject.UnixAbsoluteFilePath,
	destinationPath tkValueObject.UnixAbsoluteFilePath,
	shouldOverwrite bool,
) MoveUnixFile {
	return MoveUnixFile{
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
		ShouldOverwrite: shouldOverwrite,
	}
}
