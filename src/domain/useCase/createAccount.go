package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func CreateAccount(
	accQueryRepo repository.AccQueryRepo,
	accCmdRepo repository.AccCmdRepo,
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	createAccount dto.CreateAccount,
) error {
	_, err := accQueryRepo.GetByUsername(createAccount.Username)
	if err == nil {
		return errors.New("UsernameAlreadyExists")
	}

	err = accCmdRepo.Add(createAccount)
	if err != nil {
		return errors.New("CreateAccountError")
	}

	log.Printf("Account '%v' created.", createAccount.Username.String())

	trashPath, _ := valueObject.NewUnixFilePath("/.trash")
	_, err = filesQueryRepo.GetOne(trashPath)
	if err != nil {
		trashDirPermissions, _ := valueObject.NewUnixFilePermissions("775")
		trashDirMimeType, _ := valueObject.NewMimeType("directory")
		createTrashDir := dto.NewCreateUnixFile(
			trashPath,
			trashDirPermissions,
			trashDirMimeType,
		)

		filesCmdRepo.Create(createTrashDir)
	}

	return nil
}
