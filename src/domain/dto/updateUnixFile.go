package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFile struct {
	SourcePaths     []valueObject.UnixFilePath       `json:"sourcePaths"`
	DestinationPath *valueObject.UnixFilePath        `json:"destinationPath"`
	Permissions     *valueObject.UnixFilePermissions `json:"permissions"`
	EncodedContent  *valueObject.EncodedContent      `json:"encodedContent"`
}

func NewUpdateUnixFile(
	sourcePaths []valueObject.UnixFilePath,
	destinationPath *valueObject.UnixFilePath,
	permissions *valueObject.UnixFilePermissions,
	encodedContent *valueObject.EncodedContent,
) UpdateUnixFile {
	return UpdateUnixFile{
		SourcePaths:     sourcePaths,
		DestinationPath: destinationPath,
		Permissions:     permissions,
		EncodedContent:  encodedContent,
	}
}
