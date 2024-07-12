package filesInfra

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

	t.Run("CreateUnixDirectory", func(t *testing.T) {
		dirPermissions := valueObject.NewUnixFilePermissionsPanic("0777")

		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")

		createUnixFile := dto.NewCreateUnixFile(
			filePath,
			&dirPermissions,
			valueObject.NewMimeTypePanic("directory"),
		)

		err := filesCmdRepo.Create(createUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CreateUnixFile", func(t *testing.T) {
		filePermissions := valueObject.NewUnixFilePermissionsPanic("0777")

		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")

		createUnixFile := dto.NewCreateUnixFile(
			filePath,
			&filePermissions,
			valueObject.NewMimeTypePanic("generic"),
		)

		err := filesCmdRepo.Create(createUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFileContent", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")

		err := filesCmdRepo.UpdateContent(
			filePath,
			valueObject.NewEncodedContentPanic("Q29udGVudCB0byB0ZXN0"),
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixDirectoryPermissions", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")
		filePermissions := valueObject.NewUnixFilePermissionsPanic("0777")

		err := filesCmdRepo.UpdatePermissions(filePath, filePermissions)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFilePermissions", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		filePermissions := valueObject.NewUnixFilePermissionsPanic("0777")

		err := filesCmdRepo.UpdatePermissions(filePath, filePermissions)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixDirectory", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")
		destinationFilePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir_")

		err := filesCmdRepo.Move(sourceFilePath, destinationFilePath, true)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixFile", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/filesCmdRepoTest.txt",
		)

		err := filesCmdRepo.Move(sourceFilePath, destinationFilePath, false)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixDirectory", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir_")
		destinationFilePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")

		copyUnixFileDto := dto.NewCopyUnixFile(sourceFilePath, destinationFilePath, true)

		err := filesCmdRepo.Copy(copyUnixFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixFile", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)

		copyUnixFileDto := dto.NewCopyUnixFile(sourceFilePath, destinationFilePath, false)

		err := filesCmdRepo.Copy(copyUnixFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile (with compression type)", func(t *testing.T) {
		var compressionTypePtr *valueObject.UnixCompressionType
		compressionType := valueObject.NewUnixCompressionTypePanic("gzip")
		compressionTypePtr = &compressionType

		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirCompress",
		)

		compressUnixFiles := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{sourceFilePath},
			destinationFilePath,
			compressionTypePtr,
		)

		_, err := filesCmdRepo.Compress(compressUnixFiles)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile (without compression type)", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirCompress",
		)

		compressUnixFiles := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{sourceFilePath},
			destinationFilePath,
			nil,
		)

		_, err := filesCmdRepo.Compress(compressUnixFiles)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile (with compression type in file path)", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirCompress_.gzip",
		)

		compressUnixFiles := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{sourceFilePath},
			destinationFilePath,
			nil,
		)

		_, err := filesCmdRepo.Compress(compressUnixFiles)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ExtractUnixFile", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirCompress.tar.gz",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirExtracted",
		)

		extractFileDto := dto.NewExtractUnixFiles(sourceFilePath, destinationFilePath)

		err := filesCmdRepo.Extract(extractFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
