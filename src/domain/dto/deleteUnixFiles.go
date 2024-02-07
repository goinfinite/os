package dto

import "github.com/speedianet/os/src/domain/valueObject"

type DeleteUnixFile struct {
	SourcePaths     []valueObject.UnixFilePath `json:"sourcePaths"`
	PermanentDelete bool                       `json:"PermanentDelete"`
}

func NewDeleteUnixFile(
	sourcePaths []valueObject.UnixFilePath,
	permanentDelete bool,
) DeleteUnixFile {
	return DeleteUnixFile{
		SourcePaths:     sourcePaths,
		PermanentDelete: permanentDelete,
	}
}
