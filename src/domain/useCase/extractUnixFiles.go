package useCase

import (
	"errors"
	"log"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

func ExtractUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	extractUnixFiles dto.ExtractUnixFiles,
) error {
	err := filesCmdRepo.Extract(extractUnixFiles)
	if err != nil {
		log.Printf("ExtractUnixFilesInfraError: %s", err.Error())
		return errors.New("ExtractUnixFilesInfraError")
	}

	log.Printf(
		"File '%s' extracted to '%s'.",
		extractUnixFiles.SourcePath.String(),
		extractUnixFiles.DestinationPath.String(),
	)

	return nil
}
