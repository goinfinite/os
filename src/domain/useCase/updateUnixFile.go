package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

type UpdateUnixFile struct {
	filesCmdRepo repository.FilesCmdRepo
}

func NewUpdateUnixFile(
	filesCmdRepo repository.FilesCmdRepo,
) UpdateUnixFile {
	return UpdateUnixFile{
		filesCmdRepo: filesCmdRepo,
	}
}

func (uc UpdateUnixFile) updateUnixFilePermissions(
	updateUnixFile dto.UpdateUnixFile,
) error {
	err := uc.filesCmdRepo.UpdatePermissions(
		updateUnixFile.SourcePath,
		*updateUnixFile.Permissions,
	)
	if err != nil {
		log.Printf("UpdateUnixFilePermissionsInfraError: %s", err.Error())
		return errors.New("UpdateUnixFilePermissionsInfraError")
	}

	log.Printf(
		"File '%s' (%s) permissions updated to '%s'.",
		updateUnixFile.SourcePath.GetFileName().String(),
		updateUnixFile.SourcePath.GetFileDir().String(),
		updateUnixFile.Permissions.String(),
	)

	return nil
}

func (uc UpdateUnixFile) updateUnixFilePath(
	updateUnixFile dto.UpdateUnixFile,
) error {
	shouldOverwrite := false
	err := uc.filesCmdRepo.Move(updateUnixFile, shouldOverwrite)
	if err != nil {
		log.Printf("MoveUnixFileInfraError: %s", err.Error())
		return errors.New("MoveUnixFileInfraError")
	}

	log.Printf(
		"File '%s' moved from %s to '%s'.",
		updateUnixFile.SourcePath.GetFileName().String(),
		updateUnixFile.SourcePath.GetFileDir().String(),
		updateUnixFile.DestinationPath.GetFileDir().String(),
	)

	return nil
}

func (uc UpdateUnixFile) Execute(
	updateUnixFile dto.UpdateUnixFile,
) error {
	if updateUnixFile.Permissions != nil {
		err := uc.updateUnixFilePermissions(updateUnixFile)
		if err != nil {
			return err
		}
	}

	if updateUnixFile.DestinationPath != nil {
		err := uc.updateUnixFilePath(updateUnixFile)
		if err != nil {
			return err
		}
	}

	if updateUnixFile.EncodedContent == nil {
		return nil
	}

	err := uc.filesCmdRepo.UpdateContent(updateUnixFile)
	if err != nil {
		log.Printf("UpdateUnixFileContentInfraError: %s", err.Error())
		return errors.New("UpdateUnixFileContentInfraError")
	}

	log.Printf(
		"File '%s' content updated.",
		updateUnixFile.SourcePath.GetFileName().String(),
	)

	return nil
}
