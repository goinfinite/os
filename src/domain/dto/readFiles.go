package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadFilesRequest struct {
	SourcePath            tkValueObject.UnixAbsoluteFilePath `json:"sourcePath"`
	ShouldIncludeFileTree *bool                              `json:"shouldIncludeFileTree,omitempty"`
}

type ReadFilesResponse struct {
	FileTree *UnixFileBranch   `json:"fileTree"`
	Files    []entity.UnixFile `json:"files"`
}
