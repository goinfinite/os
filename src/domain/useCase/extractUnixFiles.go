package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func ExtractUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	extractUnixFiles dto.ExtractUnixFiles,
) error {
	unixFiles, err := filesQueryRepo.Get(extractUnixFiles.Path)
	if err != nil || len(unixFiles) < 1 {
		return errors.New("FileDoesNotExists")
	}

	unixDestinationFiles, err := filesQueryRepo.Get(extractUnixFiles.DestinationPath)
	if err == nil || len(unixDestinationFiles) > 0 {
		return errors.New("DestinationAlreadyExists")
	}

	err = filesCmdRepo.Extract(
		extractUnixFiles.Path,
		extractUnixFiles.DestinationPath,
	)
	if err != nil {
		return errors.New("UnableToExtractFilesAndDirectories")
	}

	log.Printf(
		"File '%s' extracted to '%s'.",
		extractUnixFiles.Path.String(),
		extractUnixFiles.DestinationPath.String(),
	)

	return nil
}
