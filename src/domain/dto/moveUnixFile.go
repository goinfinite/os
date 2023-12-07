package dto

import "github.com/speedianet/os/src/domain/valueObject"

type MoveUnixFile struct {
	OriginPath  valueObject.UnixFilePath `json:"originPath"`
	DestinyPath valueObject.UnixFilePath `json:"destinyPath"`
}

func NewMoveUnixFile(
	OriginPath valueObject.UnixFilePath,
	DestinyPath valueObject.UnixFilePath,
) MoveUnixFile {
	return MoveUnixFile{
		OriginPath:  OriginPath,
		DestinyPath: DestinyPath,
	}
}
