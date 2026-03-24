package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type UpdateUnixFileOwnership struct {
	SourcePath  tkValueObject.UnixAbsoluteFilePath `json:"sourcePath"`
	Ownership   tkValueObject.UnixFileOwnership    `json:"ownership"`
	IsRecursive bool                               `json:"isRecursive"`
}

func NewUpdateUnixFileOwnership(
	sourcePath tkValueObject.UnixAbsoluteFilePath,
	ownership tkValueObject.UnixFileOwnership,
	isRecursive bool,
) UpdateUnixFileOwnership {
	return UpdateUnixFileOwnership{
		SourcePath:  sourcePath,
		Ownership:   ownership,
		IsRecursive: isRecursive,
	}
}
