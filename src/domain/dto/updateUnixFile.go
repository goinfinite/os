package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFile struct {
	Path            valueObject.UnixFilePath         `json:"path"`
	DestinationPath *valueObject.UnixFilePath        `json:"destinationPath"`
	Permissions     *valueObject.UnixFilePermissions `json:"permissions"`
}

func NewUpdateUnixFile(
	path valueObject.UnixFilePath,
	destinationPath *valueObject.UnixFilePath,
	permissions *valueObject.UnixFilePermissions,
) UpdateUnixFile {
	return UpdateUnixFile{
		Path:            path,
		DestinationPath: destinationPath,
		Permissions:     permissions,
	}
}
