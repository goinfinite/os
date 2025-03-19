package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadFilesRequest struct {
	SourcePath            valueObject.UnixFilePath `json:"serviceName"`
	ShouldIncludeFileTree *bool                    `json:"shouldIncludeFileTree,omitempty"`
}

type ReadFilesResponse struct {
	FileTree UnixFileTree      `json:"fileTree"`
	Files    []entity.UnixFile `json:"files"`
}
