package infra

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestFilesCmdRepo(t *testing.T) {
	filesCmdRepo := FilesCmdRepo{}

	currentUser, _ := user.Current()
	fileBasePathStr := fmt.Sprintf("/home/%s", currentUser.Username)

	t.Run("AddUnixDirectory", func(t *testing.T) {
		addUnixFile := dto.NewAddUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
			valueObject.NewUnixFilePermissionsPanic("0777"),
			valueObject.NewUnixFileTypePanic("directory"),
		)

		err := filesCmdRepo.Add(addUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("AddUnixFile", func(t *testing.T) {
		addUnixFile := dto.NewAddUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir/filesCmdRepoTest.txt"),
			valueObject.NewUnixFilePermissionsPanic("0777"),
			valueObject.NewUnixFileTypePanic("file"),
		)

		err := filesCmdRepo.Add(addUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixDirectory", func(t *testing.T) {
		moveUnixFile := dto.NewMoveUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_"),
			valueObject.NewUnixFileTypePanic("directory"),
		)

		err := filesCmdRepo.Move(moveUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixDirectory", func(t *testing.T) {
		moveUnixFile := dto.NewMoveUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/filesCmdRepoTest.txt"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/filesCmdRepoTest_.txt"),
			valueObject.NewUnixFileTypePanic("file"),
		)

		err := filesCmdRepo.Move(moveUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFileContent", func(t *testing.T) {
		updateUnixFileContent := dto.NewUpdateUnixFileContent(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/filesCmdRepoTest_.txt"),
			valueObject.NewUnixFileContentPanic("Q29udGVudCB0byB0ZXN0"),
		)

		err := filesCmdRepo.UpdateContent(updateUnixFileContent)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixDirectoryPermissions", func(t *testing.T) {
		filePath := valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir_")
		filePermissions := valueObject.NewUnixFilePermissionsPanic("0777")
		fileType := valueObject.NewUnixFileTypePanic("directory")

		err := filesCmdRepo.UpdatePermissions(
			filePath,
			filePermissions,
			fileType,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFilePermissions", func(t *testing.T) {
		filePath := valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir_/filesCmdRepoTest_.txt")
		filePermissions := valueObject.NewUnixFilePermissionsPanic("0777")
		fileType := valueObject.NewUnixFileTypePanic("file")

		err := filesCmdRepo.UpdatePermissions(
			filePath,
			filePermissions,
			fileType,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
