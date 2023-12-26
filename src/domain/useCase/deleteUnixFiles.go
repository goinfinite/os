package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	unixFilePaths []valueObject.UnixFilePath,
) {
	for fileToDeleteIndex, fileToDelete := range unixFilePaths {
		isRootPath := fileToDelete.String() == "/"
		if !isRootPath {
			continue
		}

		log.Printf("Path '/' is not allowed to delete.")

		filesToDeleteBeforeNotAllowedPath := unixFilePaths[:fileToDeleteIndex]
		filesToDeleteAfterNotAllowedPath := unixFilePaths[fileToDeleteIndex+1:]
		filesToDeleteWithoutNotAllowedPath := append(
			filesToDeleteBeforeNotAllowedPath,
			filesToDeleteAfterNotAllowedPath...,
		)

		unixFilePaths = filesToDeleteWithoutNotAllowedPath
	}

	filesCmdRepo.Delete(unixFilePaths)
}
