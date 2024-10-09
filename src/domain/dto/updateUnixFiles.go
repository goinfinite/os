package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateUnixFiles struct {
	SourcePaths     []valueObject.UnixFilePath       `json:"sourcePaths"`
	DestinationPath *valueObject.UnixFilePath        `json:"destinationPath"`
	Permissions     *valueObject.UnixFilePermissions `json:"permissions"`
	EncodedContent  *valueObject.EncodedContent      `json:"encodedContent"`
}

func NewUpdateUnixFiles(
	sourcePaths []valueObject.UnixFilePath,
	destinationPath *valueObject.UnixFilePath,
	permissions *valueObject.UnixFilePermissions,
	encodedContent *valueObject.EncodedContent,
) UpdateUnixFiles {
	return UpdateUnixFiles{
		SourcePaths:     sourcePaths,
		DestinationPath: destinationPath,
		Permissions:     permissions,
		EncodedContent:  encodedContent,
	}
}
