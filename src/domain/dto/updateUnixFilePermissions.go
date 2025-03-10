package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateUnixFilePermissions struct {
	SourcePath           valueObject.UnixFilePath         `json:"sourcePath"`
	FilePermissions      valueObject.UnixFilePermissions  `json:"filePermissions"`
	DirectoryPermissions *valueObject.UnixFilePermissions `json:"directoryPermissions"`
}

func NewUpdateUnixFilePermissions(
	sourcePath valueObject.UnixFilePath,
	filePermissions valueObject.UnixFilePermissions,
	directoryPermissions *valueObject.UnixFilePermissions,
) UpdateUnixFilePermissions {
	return UpdateUnixFilePermissions{
		SourcePath:           sourcePath,
		FilePermissions:      filePermissions,
		DirectoryPermissions: directoryPermissions,
	}
}
