package infra

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestFilesCmdRepo(t *testing.T) {
	filesCmdRepo := NewFilesCmdRepo()

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
		err := filesCmdRepo.Move(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_"),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixDirectory", func(t *testing.T) {
		err := filesCmdRepo.Move(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/filesCmdRepoTest.txt"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/filesCmdRepoTest_.txt"),
		)
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

		err := filesCmdRepo.UpdatePermissions(
			filePath,
			filePermissions,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFilePermissions", func(t *testing.T) {
		filePath := valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir_/filesCmdRepoTest_.txt")
		filePermissions := valueObject.NewUnixFilePermissionsPanic("0777")

		err := filesCmdRepo.UpdatePermissions(
			filePath,
			filePermissions,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixDirectory", func(t *testing.T) {
		addUnixFileCopyDto := dto.NewAddUnixFileCopy(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
		)

		err := filesCmdRepo.Copy(addUnixFileCopyDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixFile", func(t *testing.T) {
		addUnixFileCopyDto := dto.NewAddUnixFileCopy(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir/testDir_/filesCmdRepoTest_.txt"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir/testDir_/filesCmdRepoTest.txt"),
		)

		err := filesCmdRepo.Copy(addUnixFileCopyDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile", func(t *testing.T) {
		err := filesCmdRepo.Compress(
			[]valueObject.UnixFilePath{valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir_/filesCmdRepoTest_.txt")},
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress"),
			valueObject.NewUnixCompressionTypePanic("gzip"),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ExtractUnixFile", func(t *testing.T) {
		err := filesCmdRepo.Extract(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress.gzip"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirExtracted"),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("DeleteUnixDirectory", func(t *testing.T) {
		err := filesCmdRepo.Delete(
			valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir/testDir_"),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("DeleteUnixFile", func(t *testing.T) {
		err := filesCmdRepo.Delete(
			valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir/filesCmdRepoTest_.txt"),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
