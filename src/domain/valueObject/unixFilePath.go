package valueObject

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

var UnixFilePathFileSystemRootDir = tkValueObject.UnixAbsoluteFilePath("/")
var UnixFilePathAppWorkingDir = tkValueObject.UnixAbsoluteFilePath("/app")
var UnixFilePathTrashDir = tkValueObject.UnixAbsoluteFilePath("/app/.trash")

func NewUnixFilePath(
	value interface{},
) (tkValueObject.UnixAbsoluteFilePath, error) {
	return tkValueObject.NewUnixAbsoluteFilePath(value, false)
}
