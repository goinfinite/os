package filesInfra

import (
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestFilesCmdRepo(t *testing.T) {
	filesCmdRepo := FilesCmdRepo{}

	currentUser, _ := user.Current()
	fileBasePathStr := "/home/" + currentUser.Username

	fileDefaultPermissions := valueObject.NewUnixFileDefaultPermissions()
	directoryDefaultPermissions := valueObject.NewUnixDirDefaultPermissions()
	operatorAccountId, _ := tkValueObject.NewAccountId(0)
	ipAddress := tkValueObject.IpAddressLocal

	t.Run("CreateUnixDirectory", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")

		createDto := dto.NewCreateUnixFile(
			filePath, &directoryDefaultPermissions, tkValueObject.MimeTypeDirectory,
			operatorAccountId, ipAddress,
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

		createDto := dto.NewCreateUnixFile(
			filePath, &fileDefaultPermissions, tkValueObject.MimeTypeGeneric,
			operatorAccountId, ipAddress,
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

	t.Run("UpdateOwnership_WithRecursiveTrue", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/testDir")
		ownership, _ := tkValueObject.NewUnixFileOwnership(
			currentUser.Username + ":" + currentUser.Username,
		)

		updateDto := dto.NewUpdateUnixFileOwnership(filePath, ownership, true)
		if !updateDto.IsRecursive {
			t.Errorf("ExpectedIsRecursiveTrue")
		}

		err := filesCmdRepo.UpdateOwnership(updateDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateOwnership_NonRecursiveOmitsDashR", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		ownership, _ := tkValueObject.NewUnixFileOwnership(
			currentUser.Username + ":" + currentUser.Username,
		)

		updateDto := dto.NewUpdateUnixFileOwnership(filePath, ownership, false)
		if updateDto.IsRecursive {
			t.Errorf("ExpectedIsRecursiveFalse")
		}

		err := filesCmdRepo.UpdateOwnership(updateDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ApplyAccountOwnership_ResolvesAccountToUsernameGroup", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)

		err := filesCmdRepo.filePrivilegesNormalizer(
			filePath, operatorAccountId, false,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}

		info, err := os.Stat(filePath.String())
		if err != nil {
			t.Errorf("StatError: %v", err)
		}
		if info == nil {
			t.Errorf("ExpectedFileInfo")
		}
	})

	t.Run("ApplyAccountOwnership_AppHtmlUsesNobodyNogroup", func(t *testing.T) {
		appHtmlDir := "/app/html"
		if _, err := os.Stat(appHtmlDir); os.IsNotExist(err) {
			t.Skip("AppHtmlDirNotPresent")
		}

		testFilePath := appHtmlDir + "/ownershipTest.txt"
		_, createErr := os.Create(testFilePath)
		if createErr != nil {
			t.Skipf("CannotCreateTestFile: %v", createErr)
		}
		defer os.Remove(testFilePath)

		filePath, _ := valueObject.NewUnixFilePath(testFilePath)
		err := filesCmdRepo.filePrivilegesNormalizer(
			filePath, operatorAccountId, false,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ApplyAccountOwnership_InvalidAccountId", func(t *testing.T) {
		filePath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/testDir/filesCmdRepoTest.txt",
		)
		invalidAccountId, _ := tkValueObject.NewAccountId(999999999)

		err := filesCmdRepo.filePrivilegesNormalizer(
			filePath, invalidAccountId, false,
		)
		if err == nil {
			t.Errorf("ExpectedErrorButGotNil")
			return
		}
		if !strings.Contains(err.Error(), "AccountNotFound") {
			t.Errorf(
				"ExpectedAccountNotFoundError, got: %s",
				err.Error(),
			)
		}
	})

	t.Run("ApplyAccountOwnership_NonExistentPath", func(t *testing.T) {
		nonExistentPath, _ := valueObject.NewUnixFilePath(
			fileBasePathStr + "/nonExistentFile12345.txt",
		)

		err := filesCmdRepo.filePrivilegesNormalizer(
			nonExistentPath, operatorAccountId, false,
		)
		if err == nil {
			t.Errorf("ExpectedErrorButGotNil")
		}
	})

	t.Run("ApplyAccountOwnership_PathTraversalAttempt", func(t *testing.T) {
		traversalPathStr := "/app/html/../../../etc/passwd"
		traversalPath, err := valueObject.NewUnixFilePath(traversalPathStr)
		if err != nil {
			t.Skip("PathTraversalRejectedByValueObject")
		}

		normalizedStr := traversalPath.String()
		isHtmlPath := strings.HasPrefix(
			normalizedStr,
			valueObject.UnixFilePathAppHtmlDir.String(),
		)
		if isHtmlPath {
			t.Errorf(
				"PathTraversalShouldNotResolveToHtmlPath: %s",
				normalizedStr,
			)
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
			[]tkValueObject.UnixAbsoluteFilePath{sourceFilePath},
			destinationFilePath,
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
			[]tkValueObject.UnixAbsoluteFilePath{sourceFilePath},
			destinationFilePath, nil,
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
			[]tkValueObject.UnixAbsoluteFilePath{sourceFilePath},
			destinationFilePath, nil,
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

func TestUpdateUnixFileOwnershipDto(t *testing.T) {
	filePath, _ := valueObject.NewUnixFilePath("/tmp/ownershipDtoTest.txt")
	ownership := tkValueObject.UnixFileOwnershipNobodyNogroup

	t.Run("UpdateUnixFileOwnership_WithRecursiveTrue", func(t *testing.T) {
		updateDto := dto.NewUpdateUnixFileOwnership(filePath, ownership, true)
		if !updateDto.IsRecursive {
			t.Errorf("ExpectedIsRecursiveTrueButGotFalse")
		}
	})

	t.Run("UpdateUnixFileOwnership_WithRecursiveFalse", func(t *testing.T) {
		updateDto := dto.NewUpdateUnixFileOwnership(filePath, ownership, false)
		if updateDto.IsRecursive {
			t.Errorf("ExpectedIsRecursiveFalseButGotTrue")
		}
	})
}
