package filesInfra

import (
	"os/user"
	"testing"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestFilesCmdRepo(t *testing.T) {
	filesCmdRepo := FilesCmdRepo{}

	currentUser, _ := user.Current()
	fileBasePathStr := "/home/" + currentUser.Username

	fileDefaultPermissions := valueObject.NewUnixFileDefaultPermissions()
	directoryDefaultPermissions := valueObject.NewUnixDirDefaultPermissions()
	operatorAccountId, _ := valueObject.NewAccountId(0)
	ipAddress := valueObject.NewLocalhostIpAddress()

	t.Run("CreateUnixDirectory", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")
		mimeType, _ := valueObject.NewMimeType("directory")

		createDto := dto.NewCreateUnixFile(
			filePath, &directoryDefaultPermissions, mimeType, operatorAccountId,
			ipAddress,
		)

		err := filesCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CreateUnixFile", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		mimeType, _ := valueObject.NewMimeType("generic")

		createDto := dto.NewCreateUnixFile(
			filePath, &fileDefaultPermissions, mimeType, operatorAccountId,
			ipAddress,
		)

		err := filesCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixFileContent", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		encodedContent, _ := valueObject.NewEncodedContent("Q29udGVudCB0byB0ZXN0")

		updateContentDto := dto.NewUpdateUnixFileContent(filePath, encodedContent)

		err := filesCmdRepo.UpdateContent(updateContentDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateOnlyUnixFilePermissions", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)

		updatePermissionsDto := dto.NewUpdateUnixFilePermissions(
			filePath, fileDefaultPermissions, nil,
		)

		err := filesCmdRepo.UpdatePermissions(updatePermissionsDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateUnixDirectoryAndFilePermissions", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")

		updatePermissionsDto := dto.NewUpdateUnixFilePermissions(
			filePath, fileDefaultPermissions, &directoryDefaultPermissions,
		)

		err := filesCmdRepo.UpdatePermissions(updatePermissionsDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixDirectory", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_",
		)

		moveDto := dto.NewMoveUnixFile(sourceFilePath, destinationFilePath, true)

		err := filesCmdRepo.Move(moveDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixFile", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/filesCmdRepoTest.txt",
		)

		moveDto := dto.NewMoveUnixFile(sourceFilePath, destinationFilePath, false)

		err := filesCmdRepo.Move(moveDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixDirectory", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir_")
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir",
		)

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
			fileBasePathStr + "/filesCmdRepoTest.txt",
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
		compressionType, _ := valueObject.NewUnixCompressionType("tgz")
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirCompressWithType",
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
			fileBasePathStr + "/testDir_/testDirCompressWithoutType",
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

	t.Run("CompressUnixFile (with compression type in file name)", func(t *testing.T) {
		sourceFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		destinationFilePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir_/testDirCompressWithTypeOnFileName_.gzip",
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
			fileBasePathStr + "/testDir_/testDirCompressWithType.tgz",
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
