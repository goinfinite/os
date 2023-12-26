package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddUnixFile struct {
	SourcePath  valueObject.UnixFilePath        `json:"sourcePath"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	Type        valueObject.UnixFileType        `json:"type"`
}

func NewAddUnixFile(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
	fileType valueObject.UnixFileType,
) AddUnixFile {
	return AddUnixFile{
		SourcePath:  sourcePath,
		Permissions: permissions,
		Type:        fileType,
	}
}
