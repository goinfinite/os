package filesInfra

import (
	"os"
	"os/user"
	"testing"

	"github.com/goinfinite/os/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func TestFilesQueryRepo(t *testing.T) {
	filesQueryRepo := NewFilesQueryRepo()
	fileClerk := tkInfra.FileClerk{}
	currentUser, _ := user.Current()
	userHomeDir := "/home/" + currentUser.Username

	// Note: Setup/teardown are intentionally inline — test independence
	// requires each file to own its preconditions, even if it duplicates code.
	homeDirCreationErr := os.MkdirAll(userHomeDir, 0755)
	if homeDirCreationErr != nil {
		t.Fatalf("HomeDirCreationFailed: %v", homeDirCreationErr)
	}
	t.Cleanup(func() {
		removeErr := os.RemoveAll(userHomeDir)
		if removeErr != nil {
			t.Errorf("HomeDirCleanupFailed: %v", removeErr)
		}
	})

	t.Run("Read", func(t *testing.T) {
		unixDirPath, _ := tkValueObject.NewUnixAbsoluteFilePath(userHomeDir, false)
		requestDto := dto.ReadFilesRequest{
			SourcePath: unixDirPath,
		}

		_, err := filesQueryRepo.Read(requestDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
	})

	t.Run("ReadWithInvalidDirectory", func(t *testing.T) {
		invalidUnixPath, _ := tkValueObject.NewUnixAbsoluteFilePath("/aaa/bbb/ccc", false)
		requestDto := dto.ReadFilesRequest{
			SourcePath: invalidUnixPath,
		}

		_, err := filesQueryRepo.Read(requestDto)
		if err == nil {
			t.Errorf("ExpectedErrorButGotNil")
		}
	})

	t.Run("SimplifiedRootBranchName", func(t *testing.T) {
		rootPath, _ := tkValueObject.NewUnixAbsoluteFilePath("/", false)
		simplifiedRoot, err := filesQueryRepo.simplifiedUnixFileFactory(rootPath)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		if simplifiedRoot.Name.String() != "root" {
			t.Errorf("ExpectedRootNameButGot: %s", simplifiedRoot.Name.String())
		}
	})

	t.Run("ReadFollowingSymlink", func(t *testing.T) {
		downloadsDirPathStr := userHomeDir + "/Downloads"
		downloadsDirCreationErr := os.MkdirAll(downloadsDirPathStr, 0755)
		if downloadsDirCreationErr != nil {
			t.Fatalf("DownloadsDirCreationFailed: %v", downloadsDirCreationErr)
		}
		downloadFilePathStr := downloadsDirPathStr + "/download.txt"
		downloadFileWriteErr := os.WriteFile(
			downloadFilePathStr, []byte("download"), 0644,
		)
		if downloadFileWriteErr != nil {
			t.Fatalf("DownloadFileCreationFailed: %v", downloadFileWriteErr)
		}
		t.Cleanup(func() {
			removeErr := os.RemoveAll(downloadsDirPathStr)
			if removeErr != nil {
				t.Errorf("DownloadsDirCleanupFailed: %v", removeErr)
			}
		})

		downloadsDirPath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			downloadsDirPathStr, false,
		)
		tmpSymlinkPath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			userHomeDir+"/tmpSymlink", false,
		)
		requestDto := dto.ReadFilesRequest{
			SourcePath: tmpSymlinkPath,
		}

		symlinkCreationErr := fileClerk.CreateSymlink(
			downloadsDirPath.String(), tmpSymlinkPath.String(), false,
		)
		if symlinkCreationErr != nil {
			t.Fatalf("ExpectedNoErrorButGot: %s", symlinkCreationErr.Error())
		}
		t.Cleanup(func() {
			removeErr := fileClerk.RemoveSymlink(tmpSymlinkPath.String())
			if removeErr != nil {
				t.Errorf("SymlinkCleanupFailed: %v", removeErr)
			}
		})

		responseDto, err := filesQueryRepo.Read(requestDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		if len(responseDto.Files) == 0 {
			t.Errorf("ExpectedNonEmptyFilesButGotEmpty")
		}
	})

	t.Run("ReadFirstFile", func(t *testing.T) {
		unixFilePathStr := userHomeDir + "/filesQueryRepoTest.txt"
		fileWriteErr := os.WriteFile(
			unixFilePathStr, []byte("query test"), 0644,
		)
		if fileWriteErr != nil {
			t.Fatalf("TestFileCreationFailed: %v", fileWriteErr)
		}
		t.Cleanup(func() {
			removeErr := os.Remove(unixFilePathStr)
			if removeErr != nil {
				t.Errorf("TestFileCleanupFailed: %v", removeErr)
			}
		})

		unixFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			unixFilePathStr, false,
		)
		_, err := filesQueryRepo.ReadFirst(unixFilePath)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
	})
}
