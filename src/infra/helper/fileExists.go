package infraHelper

import tkInfra "github.com/goinfinite/tk/src/infra"

var fileClerk = tkInfra.FileClerk{}

func FileExists(filePath string) bool {
	return fileClerk.FileExists(filePath)
}
