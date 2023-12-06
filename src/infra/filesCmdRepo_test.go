package infra

import (
	"testing"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestFilesCmdRepo(t *testing.T) {
	filesCmdRepo := FilesCmdRepo{}

	t.Run("AddUnixDirectory", func(t *testing.T) {
		addUnixFile := dto.NewAddUnixFile(
			valueObject.NewUnixFilePathPanic("/home/mmp/testDir"),
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
			valueObject.NewUnixFilePathPanic("/home/mmp/testDir/filesCmdRepoTest.txt"),
			valueObject.NewUnixFilePermissionsPanic("0777"),
			valueObject.NewUnixFileTypePanic("file"),
		)

		err := filesCmdRepo.Add(addUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixDirectoryPermissions", func(t *testing.T) {
		filePath := valueObject.NewUnixFilePathPanic("/home/mmp/testDir")
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
		filePath := valueObject.NewUnixFilePathPanic("/home/mmp/testDir/filesCmdRepoTest.txt")
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
