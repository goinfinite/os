package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type FixUnixFilePermissions struct {
	SourcePath           valueObject.UnixFilePath        `json:"sourcePath"`
	DirectoryPermissions valueObject.UnixFilePermissions `json:"directoryPermissions"`
	FilePermissions      valueObject.UnixFilePermissions `json:"filePermissions"`
}

func NewFixUnixFilePermissions(
	sourcePath valueObject.UnixFilePath,
	directoryPermissions valueObject.UnixFilePermissions,
	filePermissions valueObject.UnixFilePermissions,
) FixUnixFilePermissions {
	return FixUnixFilePermissions{
		SourcePath:           sourcePath,
		DirectoryPermissions: directoryPermissions,
		FilePermissions:      filePermissions,
	}
}
