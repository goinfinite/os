package dto

import "github.com/speedianet/os/src/domain/valueObject"

type DeleteUnixFiles struct {
	SourcePaths []valueObject.UnixFilePath `json:"sourcePaths"`
	HardDelete  bool                       `json:"hardDelete"`
}

func NewDeleteUnixFiles(
	sourcePaths []valueObject.UnixFilePath,
	hardDelete bool,
) DeleteUnixFiles {
	return DeleteUnixFiles{
		SourcePaths: sourcePaths,
		HardDelete:  hardDelete,
	}
}
