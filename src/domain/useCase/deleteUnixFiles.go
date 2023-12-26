package useCase

import (
	"log"
	"slices"

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

		log.Printf("DeleteUnixFilesError: Path '/' is not allowed to delete.")

		fileToDeleteAfterNotAllowedPathIndex := fileToDeleteIndex + 1
		filesToDeleteWithoutNotAllowedPath := slices.Delete(
			unixFilePaths,
			fileToDeleteIndex,
			fileToDeleteAfterNotAllowedPathIndex,
		)

		unixFilePaths = filesToDeleteWithoutNotAllowedPath
	}

	filesCmdRepo.Delete(unixFilePaths)
}
