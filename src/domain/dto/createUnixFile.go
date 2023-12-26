package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateUnixFile struct {
	SourcePath  valueObject.UnixFilePath        `json:"sourcePath"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	Type        valueObject.UnixFileType        `json:"type"`
}

func NewCreateUnixFile(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
	fileType valueObject.UnixFileType,
) CreateUnixFile {
	return CreateUnixFile{
		SourcePath:  sourcePath,
		Permissions: permissions,
		Type:        fileType,
	}
}
