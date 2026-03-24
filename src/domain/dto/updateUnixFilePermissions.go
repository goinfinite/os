package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateUnixFilePermissions struct {
	SourcePath           tkValueObject.UnixAbsoluteFilePath `json:"sourcePath"`
	FilePermissions      valueObject.UnixFilePermissions    `json:"filePermissions"`
	DirectoryPermissions *valueObject.UnixFilePermissions   `json:"directoryPermissions"`
}

func NewUpdateUnixFilePermissions(
	sourcePath tkValueObject.UnixAbsoluteFilePath,
	filePermissions valueObject.UnixFilePermissions,
	directoryPermissions *valueObject.UnixFilePermissions,
) UpdateUnixFilePermissions {
	return UpdateUnixFilePermissions{
		SourcePath:           sourcePath,
		FilePermissions:      filePermissions,
		DirectoryPermissions: directoryPermissions,
	}
}
