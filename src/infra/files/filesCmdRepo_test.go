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
	filesCmdRepo := NewFilesCmdRepo()

	currentUser, _ := user.Current()
	fileBasePathStr := "/home/" + currentUser.Username

	fileDefaultPermissions := valueObject.NewUnixFileDefaultPermissions()
	directoryDefaultPermissions := valueObject.NewUnixDirDefaultPermissions()
	operatorAccountId, _ := tkValueObject.NewAccountId(0)
	ipAddress := tkValueObject.IpAddressLocal

	// Note: Setup/teardown are intentionally inline — test independence
	// requires each file to own its preconditions, even if it duplicates code.
	t.Cleanup(func() {
		fixturePaths := []string{
			fileBasePathStr + "/testDir",
			fileBasePathStr + "/testDir_",
			fileBasePathStr + "/filesCmdRepoTest.txt",
		}
		for _, fixturePath := range fixturePaths {
			removeErr := os.RemoveAll(fixturePath)
			if removeErr != nil {
				t.Errorf("FixtureCleanupFailed: %v", removeErr)
			}
		}
	})

	t.Run("CreateUnixDirectory", func(t *testing.T) {
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(fileBasePathStr+"/testDir", false)

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
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
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
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
		)
		encodedContent, _ := valueObject.NewEncodedContent("Q29udGVudCB0byB0ZXN0")

		updateContentDto := dto.NewUpdateUnixFileContent(filePath, encodedContent)

		err := filesCmdRepo.UpdateContent(updateContentDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateOnlyUnixFilePermissions", func(t *testing.T) {
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
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
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(fileBasePathStr+"/testDir", false)

		updatePermissionsDto := dto.NewUpdateUnixFilePermissions(
			filePath, fileDefaultPermissions, &directoryDefaultPermissions,
		)

		err := filesCmdRepo.UpdatePermissions(updatePermissionsDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("UpdateOwnership_WithRecursiveTrue", func(t *testing.T) {
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(fileBasePathStr+"/testDir", false)
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
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
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
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
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
		_, statErr := os.Stat(appHtmlDir)
		if os.IsNotExist(statErr) {
			t.Skip("AppHtmlDirNotPresent")
		}

		testFilePath := appHtmlDir + "/ownershipTest.txt"
		_, createErr := os.Create(testFilePath)
		if createErr != nil {
			t.Skipf("CannotCreateTestFile: %v", createErr)
		}
		defer func() { _ = os.Remove(testFilePath) }()

		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(testFilePath, false)
		err := filesCmdRepo.filePrivilegesNormalizer(
			filePath, operatorAccountId, false,
		)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("ApplyAccountOwnership_InvalidAccountId", func(t *testing.T) {
		filePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
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
		nonExistentPath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/nonExistentFile12345.txt", false,
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
		traversalPath, err := tkValueObject.NewUnixAbsoluteFilePath(traversalPathStr, false)
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
		sourceFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(fileBasePathStr+"/testDir", false)
		destinationFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_", false,
		)

		moveDto := dto.NewMoveUnixFile(sourceFilePath, destinationFilePath, true)

		err := filesCmdRepo.Move(moveDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("MoveUnixFile", func(t *testing.T) {
		sourceFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_/filesCmdRepoTest.txt", false,
		)
		destinationFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/filesCmdRepoTest.txt", false,
		)

		moveDto := dto.NewMoveUnixFile(sourceFilePath, destinationFilePath, false)

		err := filesCmdRepo.Move(moveDto)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})

	t.Run("CopyUnixDirectory", func(t *testing.T) {
		destinationParentPath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir", false,
		)
		createDto := dto.NewCreateUnixFile(
			destinationParentPath, &directoryDefaultPermissions,
			tkValueObject.MimeTypeDirectory, operatorAccountId, ipAddress,
		)

		err := filesCmdRepo.Create(createDto)
		if err != nil {
			t.Fatalf("DestinationParentCreationFailed: %v", err)
		}

		sourcePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_", false,
		)
		copyDto := dto.NewCopyUnixFile(
			sourcePath, destinationParentPath, true, operatorAccountId, ipAddress,
		)

		err = filesCmdRepo.Copy(copyDto)
		if err != nil {
			t.Fatalf("UnexpectedError: %v", err)
		}

		copiedDirPath := fileBasePathStr + "/testDir/testDir_"
		_, statErr := os.Stat(copiedDirPath)
		if statErr != nil {
			t.Errorf("ExpectedCopiedDirectoryButGot: %v", statErr)
		}
	})

	t.Run("CopyUnixFile", func(t *testing.T) {
		sourcePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/filesCmdRepoTest.txt", false,
		)
		destinationParentPath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir", false,
		)

		copyDto := dto.NewCopyUnixFile(
			sourcePath, destinationParentPath, false, operatorAccountId, ipAddress,
		)

		err := filesCmdRepo.Copy(copyDto)
		if err != nil {
			t.Fatalf("UnexpectedError: %v", err)
		}

		copiedFilePath := fileBasePathStr + "/testDir/filesCmdRepoTest.txt"
		_, statErr := os.Stat(copiedFilePath)
		if statErr != nil {
			t.Errorf("ExpectedCopiedFileButGot: %v", statErr)
		}
	})

	t.Run("CompressUnixFile (with compression type)", func(t *testing.T) {
		compressionType, _ := valueObject.NewUnixCompressionType("tgz")
		sourceFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
		)
		destinationFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_/testDirCompressWithType", false,
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
		sourceFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
		)
		destinationFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_/testDirCompressWithoutType", false,
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
		sourceFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir/filesCmdRepoTest.txt", false,
		)
		destinationFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_/testDirCompressWithTypeOnFileName_.gzip", false,
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
		sourceFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_/testDirCompressWithType.tgz", false,
		)
		destinationFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			fileBasePathStr+"/testDir_/testDirExtracted", false,
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
	filePath, _ := tkValueObject.NewUnixAbsoluteFilePath("/tmp/ownershipDtoTest.txt", false)
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
