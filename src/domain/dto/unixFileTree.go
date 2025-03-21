package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type UnixFileBranch struct {
	entity.SimplifiedUnixFile
	Branches map[valueObject.UnixFileName]UnixFileBranch `json:"branches"`
}

func NewUnixFileBranch(trunkBranchFile entity.SimplifiedUnixFile) UnixFileBranch {
	return UnixFileBranch{
		SimplifiedUnixFile: trunkBranchFile,
		Branches:           map[valueObject.UnixFileName]UnixFileBranch{},
	}
}
