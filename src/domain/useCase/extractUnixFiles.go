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
	unixFileToExtractExists, err := filesQueryRepo.Exists(extractUnixFiles.Path)
	if err != nil {
		return err
	}

	if !unixFileToExtractExists {
		return errors.New("PathDoesNotExists")
	}

	unixDestinationFileExists, err := filesQueryRepo.Exists(extractUnixFiles.DestinationPath)
	if err != nil {
		return err
	}

	if !unixDestinationFileExists {
		return errors.New("DestinationPathAlreadyExists")
	}

	err = filesCmdRepo.Extract(
		extractUnixFiles.Path,
		extractUnixFiles.DestinationPath,
	)
	if err != nil {
		return err
	}

	log.Printf(
		"File '%s' extracted to '%s'.",
		extractUnixFiles.Path.String(),
		extractUnixFiles.DestinationPath.String(),
	)

	return nil
}
