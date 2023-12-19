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
		if err != nil || !unixFileExists {
			continue
		}

		err = filesCmdRepo.Delete(filePath)
		if err != nil {
			continue
		}

		isDir, err := filesQueryRepo.IsDir(filePath)
		if err != nil {
			continue
		}

		inodeName := "File"
		if isDir {
			inodeName = "Directory"
		}

		log.Printf("%s '%s' deleted.", inodeName, filePath.String())
	}
}
