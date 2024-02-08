package dto

import "github.com/speedianet/os/src/domain/valueObject"

type DeleteUnixFiles struct {
	SourcePaths     []valueObject.UnixFilePath `json:"sourcePaths"`
	PermanentDelete bool                       `json:"permanentDelete"`
}

func NewDeleteUnixFiles(
	sourcePaths []valueObject.UnixFilePath,
	permanentDelete bool,
) DeleteUnixFiles {
	return DeleteUnixFiles{
		SourcePaths:     sourcePaths,
		PermanentDelete: permanentDelete,
	}
}
