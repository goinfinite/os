package infra

import (
	"errors"
	"fmt"
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
		fileType := "File"
		fileIsDir, _ := moveUnixFile.OriginPath.IsDir()
		if fileIsDir {
			fileType = "Directory"
		}

		moveErrorStr := fmt.Sprintf("MoveUnix%sError", fileType)

		log.Printf("%s: %s", moveErrorStr, err)
		return errors.New(moveErrorStr)
	}

	return nil
}

func (repo FilesCmdRepo) UpdateContent(
	updateUnixFileContent dto.UpdateUnixFileContent,
) error {
	file, err := os.OpenFile(updateUnixFileContent.Path.String(), os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("OpenFileError: %s", err)
		return errors.New("OpenFileError")
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Printf("TruncateFileError: %s", err)
		return errors.New("TruncateFileError")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("SeekFileError: %s", err)
		return errors.New("SeekFileError")
	}

	_, err = file.WriteString(updateUnixFileContent.Content.GetDecodedContent())
	if err != nil {
		log.Printf("WriteFileError: %s", err)
		return errors.New("WriteFileError")
	}

	err = file.Sync()
	if err != nil {
		log.Printf("FileSyncError: %s", err)
		return errors.New("FileSyncError")
	}

	return nil
}

func (repo FilesCmdRepo) UpdatePermissions(
	unixFilePath valueObject.UnixFilePath,
	unixFilePermissions valueObject.UnixFilePermissions,
) error {
	err := os.Chmod(unixFilePath.String(), unixFilePermissions.GetFileMode())
	if err != nil {
		fileType := "File"
		fileIsDir, _ := unixFilePath.IsDir()
		if fileIsDir {
			fileType = "Directory"
		}

		chmodErrorStr := fmt.Sprintf("ChmodUnix%sError", fileType)

		log.Printf("%s: %s", chmodErrorStr, err)
		return errors.New(chmodErrorStr)
	}

	return nil
}
