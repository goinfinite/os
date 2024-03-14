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

		createUnixFile := dto.NewCreateUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
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

		createUnixFile := dto.NewCreateUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir/filesCmdRepoTest.txt"),
			&filePermissions,
			valueObject.NewMimeTypePanic("generic"),
		)

		err := filesCmdRepo.Create(createUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFileContent", func(t *testing.T) {
		err := filesCmdRepo.UpdateContent(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir/filesCmdRepoTest.txt"),
			valueObject.NewEncodedContentPanic("Q29udGVudCB0byB0ZXN0"),
		)
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
		err := filesCmdRepo.Move(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_"),
			true,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixFile", func(t *testing.T) {
		err := filesCmdRepo.Move(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/filesCmdRepoTest.txt"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/filesCmdRepoTest.txt"),
			false,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixDirectory", func(t *testing.T) {
		copyUnixFileDto := dto.NewCopyUnixFile(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir"),
			true,
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
			true,
		)

		err := filesCmdRepo.Copy(copyUnixFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile (with compression type)", func(t *testing.T) {
		var compressionTypePtr *valueObject.UnixCompressionType
		compressionType := valueObject.NewUnixCompressionTypePanic("gzip")
		compressionTypePtr = &compressionType

		compressUnixFiles := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")},
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress"),
			compressionTypePtr,
		)

		_, err := filesCmdRepo.Compress(compressUnixFiles)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile (without compression type)", func(t *testing.T) {
		compressUnixFiles := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")},
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress_"),
			nil,
		)

		_, err := filesCmdRepo.Compress(compressUnixFiles)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile (with compression type in file path)", func(t *testing.T) {
		compressUnixFiles := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{valueObject.NewUnixFilePathPanic(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")},
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress_.gzip"),
			nil,
		)

		_, err := filesCmdRepo.Compress(compressUnixFiles)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ExtractUnixFile", func(t *testing.T) {
		extractFileDto := dto.NewExtractUnixFiles(
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirCompress.tar.gz"),
			valueObject.NewUnixFilePathPanic(fileBasePathStr+"/testDir_/testDirExtracted"),
		)

		err := filesCmdRepo.Extract(extractFileDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
