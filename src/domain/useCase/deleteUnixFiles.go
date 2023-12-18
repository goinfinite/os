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
	for _, filePath := range unixFilePaths {
		unixFileExists, err := filesQueryRepo.Exists(filePath)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		if !unixFileExists {
			log.Printf("PathDoesNotExists: %s", filePath.String())
			continue
		}

		isDir, err := filesQueryRepo.IsDir(filePath)
		if err != nil {
			log.Printf("PathIsDirError: %v", err)
			continue
		}

		inodeName := "File"
		if isDir {
			inodeName = "Directory"
		}

		err = filesCmdRepo.Delete(filePath)
		if err != nil {
			log.Printf("Delete%sError: %v", inodeName, err)
			continue
		}

		log.Printf("%s '%s' deleted.", inodeName, filePath.String())
	}
}
