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

		err := filesCmdRepo.Create(addUnixFile)
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

		err := filesCmdRepo.Create(addUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFileContent", func(t *testing.T) {
		updateUnixFileContent := dto.NewUpdateUnixFileContent(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir/filesCmdRepoTest.txt"),
			valueObject.NewEncodedContentPanic("Q29udGVudCB0byB0ZXN0"),
		)

		err := filesCmdRepo.UpdateContent(updateUnixFileContent)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixDirectoryPermissions", func(t *testing.T) {
		filePath := valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir")
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
		filePath := valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")
		filePermissions := valueObject.NewUnixFilePermissionsPanic("0777")

		err := filesCmdRepo.UpdatePermissions(
			filePath,
			filePermissions,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixDirectory", func(t *testing.T) {
		destinationPath := valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir_")
		destinationPathPtr := &destinationPath

		permissions := valueObject.NewUnixFilePermissionsPanic("0777")
		permissionsPtr := &permissions

		updateUnixFile := dto.NewUpdateUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
			destinationPathPtr,
			permissionsPtr,
		)

		err := filesCmdRepo.Move(updateUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixFile", func(t *testing.T) {
		destinationPath := valueObject.NewUnixFilePathPanic(fileBasePathStr + "/filesCmdRepoTest.txt")
		destinationPathPtr := &destinationPath

		permissions := valueObject.NewUnixFilePermissionsPanic("0777")
		permissionsPtr := &permissions

		updateUnixFile := dto.NewUpdateUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/filesCmdRepoTest.txt"),
			destinationPathPtr,
			permissionsPtr,
		)

		err := filesCmdRepo.Move(updateUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixDirectory", func(t *testing.T) {
		copyUnixFileDto := dto.NewCopyUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
		)

		err := filesCmdRepo.Copy(copyUnixFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixFile", func(t *testing.T) {
		copyUnixFileDto := dto.NewCopyUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/filesCmdRepoTest.txt"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir/filesCmdRepoTest.txt"),
		)

		err := filesCmdRepo.Copy(copyUnixFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile", func(t *testing.T) {
		compressUnixFiles := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")},
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress"),
			valueObject.NewUnixCompressionTypePanic("gzip"),
		)

		compressionProcessReport := filesCmdRepo.Compress(compressUnixFiles)
		if len(compressionProcessReport.Failure) > 0 {
			t.Errorf("UnexpectedError: %v", compressionProcessReport.Failure[0].Reason)
		}
	})

	t.Run("ExtractUnixFile", func(t *testing.T) {
		extractFileDto := dto.NewExtractUnixFiles(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress.gzip"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirExtracted"),
		)

		err := filesCmdRepo.Extract(extractFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
