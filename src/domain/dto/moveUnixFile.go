package dto

import "github.com/speedianet/os/src/domain/valueObject"

type MoveUnixFile struct {
	OriginPath  valueObject.UnixFilePath `json:"originPath"`
	DestinyPath valueObject.UnixFilePath `json:"destinyPath"`
	Type        valueObject.UnixFileType `json:"type"`
}

func NewMoveUnixFile(
	OriginPath valueObject.UnixFilePath,
	DestinyPath valueObject.UnixFilePath,
	Type valueObject.UnixFileType,
) MoveUnixFile {
	return MoveUnixFile{
		OriginPath:  OriginPath,
		DestinyPath: DestinyPath,
		Type:        Type,
	}
}
