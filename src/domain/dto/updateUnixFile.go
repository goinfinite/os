package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFile struct {
	Path            valueObject.UnixFilePath         `json:"path"`
	DestinationPath *valueObject.UnixFilePath        `json:"destinationPath"`
	Permissions     *valueObject.UnixFilePermissions `json:"permissions"`
}

func NewUpdateUnixFile(
	Path valueObject.UnixFilePath,
	DestinationPath *valueObject.UnixFilePath,
	Permissions *valueObject.UnixFilePermissions,
) UpdateUnixFile {
	return UpdateUnixFile{
		Path:            Path,
		DestinationPath: DestinationPath,
		Permissions:     Permissions,
	}
}
