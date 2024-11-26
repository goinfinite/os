package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateUnixFilePermissions struct {
	SourcePath  valueObject.UnixFilePath        `json:"sourcePath"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
}

func NewUpdateUnixFilePermissions(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
) UpdateUnixFilePermissions {
	return UpdateUnixFilePermissions{
		SourcePath:  sourcePath,
		Permissions: permissions,
	}
}
