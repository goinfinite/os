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
	err := filesCmdRepo.Extract(extractUnixFiles)
	if err != nil {
		log.Printf("ExtractUnixFilesError: %s", err.Error())
		return errors.New("ExtractUnixFilesError")
	}

	log.Printf(
		"File '%s' extracted to '%s'.",
		extractUnixFiles.Path.String(),
		extractUnixFiles.DestinationPath.String(),
	)

	return nil
}
