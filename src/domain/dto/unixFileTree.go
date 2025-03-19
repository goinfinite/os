package dto

import "github.com/goinfinite/os/src/domain/entity"

type UnixFileTree struct {
	entity.SimplifiedUnixFile
	Children []UnixFileTree `json:"children"`
}

func NewUnixFileTree(
	simplifiedUnixFile entity.SimplifiedUnixFile,
	children []UnixFileTree,
) UnixFileTree {
	return UnixFileTree{
		SimplifiedUnixFile: simplifiedUnixFile,
		Children:           children,
	}
}

func (dto *UnixFileTree) AddUnixFile(child entity.SimplifiedUnixFile) {
	childNode := NewUnixFileTree(child, []UnixFileTree{})
	dto.Children = append(dto.Children, childNode)
}

func (dto *UnixFileTree) AddSubTree(subTree UnixFileTree) {
	dto.Children = append(dto.Children, subTree)
}
