package filesInfra

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestFilesCmdRepo(t *testing.T) {
	filesCmdRepo := FilesCmdRepo{}

	currentUser, _ := user.Current()
	fileBasePathStr := fmt.Sprintf("/home/%s", currentUser.Username)

	filePermissions, _ := valueObject.NewUnixFilePermissions("0777")
	operatorAccountId, _ := valueObject.NewAccountId(0)
	ipAddress := valueObject.NewLocalhostIpAddress()

	t.Run("CreateUnixDirectory", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")
		mimeType, _ := valueObject.NewMimeType("directory")

		createDto := dto.NewCreateUnixFile(
			filePath, &filePermissions, mimeType, operatorAccountId, ipAddress,
		)

		err := filesCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CreateUnixFile", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")
		mimeType, _ := valueObject.NewMimeType("generic")

		createDto := dto.NewCreateUnixFile(
			filePath, &filePermissions, mimeType, operatorAccountId, ipAddress,
		)

		err := filesCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFileContent", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir/filesCmdRepoTest.txt")
		encodedContent, _ := valueObject.NewEncodedContent("Q29udGVudCB0byB0ZXN0")

		err := filesCmdRepo.UpdateContent(filePath, encodedContent)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixDirectoryPermissions", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")
		filePermissions, _ := valueObject.NewUnixFilePermissions("0777")

		err := filesCmdRepo.UpdatePermissions(filePath, filePermissions)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFilePermissions", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)

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

		dto := dto.NewCopyUnixFile(
			sourceFilePath, destinationFilePath, true, operatorAccountId, ipAddress,
		)

		err := filesCmdRepo.Copy(dto)
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

		dto := dto.NewCopyUnixFile(
			sourceFilePath, destinationFilePath, false, operatorAccountId, ipAddress,
		)

		err := filesCmdRepo.Copy(dto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CompressUnixFile (with compression type)", func(t *testing.T) {
		compressionType, _ := valueObject.NewUnixCompressionType("gzip")
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirCompress",
		)

		dto := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{sourceFilePath}, destinationFilePath,
			&compressionType, operatorAccountId, ipAddress,
		)

		_, err := filesCmdRepo.Compress(dto)
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

		dto := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{sourceFilePath}, destinationFilePath, nil,
			operatorAccountId, ipAddress,
		)

		_, err := filesCmdRepo.Compress(dto)
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

		dto := dto.NewCompressUnixFiles(
			[]valueObject.UnixFilePath{sourceFilePath}, destinationFilePath, nil,
			operatorAccountId, ipAddress,
		)

		_, err := filesCmdRepo.Compress(dto)
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

		dto := dto.NewExtractUnixFiles(
			sourceFilePath, destinationFilePath, operatorAccountId, ipAddress,
		)

		err := filesCmdRepo.Extract(dto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
