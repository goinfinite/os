package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFile struct {
	SourcePath      valueObject.UnixFilePath         `json:"sourcePath"`
	DestinationPath *valueObject.UnixFilePath        `json:"destinationPath"`
	Permissions     *valueObject.UnixFilePermissions `json:"permissions"`
}

func NewUpdateUnixFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath *valueObject.UnixFilePath,
	permissions *valueObject.UnixFilePermissions,
) UpdateUnixFile {
	return UpdateUnixFile{
		SourcePath:      sourcePath,
		DestinationPath: destinationPath,
		Permissions:     permissions,
	}
}
