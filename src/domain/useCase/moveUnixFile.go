package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func MoveUnixFile(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	moveUnixFile dto.MoveUnixFile,
) error {
	fileType := "File"
	fileIsDir, _ := moveUnixFile.OriginPath.IsDir()
	if fileIsDir {
		fileType = "Directory"
	}

	unixFiles, err := filesQueryRepo.Get(moveUnixFile.OriginPath)

	if err != nil && len(unixFiles) < 1 {
		return errors.New(fileType + "DoesNotExists")
	}

	err = filesCmdRepo.Move(moveUnixFile)
	if err != nil {
		return errors.New("Move" + fileType + "Error")
	}

	fileName, _ := moveUnixFile.OriginPath.GetFileName()
	log.Printf(
		"%s '%s' moved from %s to '%s'.",
		fileType,
		fileName.String(),
		moveUnixFile.OriginPath.String(),
		moveUnixFile.DestinyPath.String(),
	)

	return nil
}
