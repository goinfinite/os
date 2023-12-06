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
	unixFiles, err := filesQueryRepo.Get(moveUnixFile.OriginPath)

	inodeType := "File"
	if moveUnixFile.Type.IsDir() {
		inodeType = "Directory"
	}

	if err != nil && len(unixFiles) < 1 {
		return errors.New(inodeType + "DoesNotExists")
	}

	err = filesCmdRepo.Move(moveUnixFile)
	if err != nil {
		return errors.New("Move" + inodeType + "Error")
	}

	fileName, _ := moveUnixFile.OriginPath.GetFileName()
	log.Printf(
		"%s '%s' moved from %s to '%s'.",
		inodeType,
		fileName.String(),
		moveUnixFile.OriginPath.String(),
		moveUnixFile.DestinyPath.String(),
	)

	return nil
}
