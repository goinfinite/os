package infra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesCmdRepo struct {
}

func (repo FilesCmdRepo) Add(addUnixFile dto.AddUnixFile) error {
	if !addUnixFile.Type.IsDir() {
		_, err := os.Create(addUnixFile.Path.String())
		if err != nil {
			log.Printf("CreateUnixFileError: %s", err)
			return errors.New("CreateUnixFileError")
		}

		return repo.UpdatePermissions(
			addUnixFile.Path,
			addUnixFile.Permissions,
			addUnixFile.Type,
		)
	}

	err := os.MkdirAll(addUnixFile.Path.String(), addUnixFile.Permissions.GetFileMode())
	if err != nil {
		log.Printf("CreateUnixFileError: %s", err)
		return errors.New("CreateUnixFileError")
	}

	return nil
}

func (repo FilesCmdRepo) Move(moveUnixFile dto.MoveUnixFile) error {
	err := os.Rename(
		moveUnixFile.OriginPath.String(),
		moveUnixFile.DestinyPath.String(),
	)
	if err != nil {
		moveErrorStr := "MoveUnixFileError"
		if moveUnixFile.Type.IsDir() {
			moveErrorStr = "MoveUnixDirectoryError"
		}

		log.Printf("%s: %s", moveErrorStr, err)
		return errors.New(moveErrorStr)
	}

	return nil
}

func (repo FilesCmdRepo) UpdatePermissions(
	unixFilePath valueObject.UnixFilePath,
	unixFilePermissions valueObject.UnixFilePermissions,
	unixFileType valueObject.UnixFileType,
) error {
	err := os.Chmod(unixFilePath.String(), unixFilePermissions.GetFileMode())
	if err != nil {
		chmodErrorStr := "ChmodUnixFileError"
		if unixFileType.IsDir() {
			chmodErrorStr = "ChmodUnixDirectoryError"
		}

		log.Printf("%s: %s", chmodErrorStr, err)
		return errors.New(chmodErrorStr)
	}

	return nil
}
