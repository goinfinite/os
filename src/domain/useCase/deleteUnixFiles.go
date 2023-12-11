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
) error {
	for _, filePath := range unixFilePaths {
		unixFiles, err := filesQueryRepo.Get(filePath)

		if err != nil || len(unixFiles) < 1 {
			log.Printf("PathDoesNotExists: %s", err)
		}

		isDir, err := filePath.IsDir()
		if err != nil {
			log.Printf("PathIsDirError: %s", err)
		}

		inodeName := "File"
		if isDir {
			inodeName = "Directory"
		}

		err = filesCmdRepo.Delete(filePath)
		if err != nil {
			log.Printf("Delete%sError: %s", inodeName, err)
		}

		log.Printf("%s '%s' deleted.", inodeName, filePath.String())
	}

	return nil
}
