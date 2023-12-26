package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func CopyUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	copyUnixFile dto.CopyUnixFile,
) error {
	err := filesCmdRepo.Copy(copyUnixFile)
	if err != nil {
		return err
	}

	fileOriginPath := copyUnixFile.OriginPath
	fileDestinationPath := copyUnixFile.DestinationPath
	log.Printf(
		"File '%s' (%s) copy added to '%s' with name '%s'.",
		fileOriginPath.GetFileName().String(),
		fileOriginPath.GetFileDir().String(),
		fileDestinationPath.GetFileName().String(),
		fileDestinationPath.GetFileDir().String(),
	)

	return nil
}
