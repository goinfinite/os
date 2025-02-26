package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateUnixFileOwnership struct {
	SourcePath valueObject.UnixFilePath      `json:"sourcePath"`
	Ownership  valueObject.UnixFileOwnership `json:"ownership"`
}

func NewUpdateUnixFileOwnership(
	sourcePath valueObject.UnixFilePath,
	ownership valueObject.UnixFileOwnership,
) UpdateUnixFileOwnership {
	return UpdateUnixFileOwnership{
		SourcePath: sourcePath,
		Ownership:  ownership,
	}
}
