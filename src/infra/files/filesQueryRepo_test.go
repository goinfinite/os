package filesInfra

import (
	"os/user"
	"testing"

	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

func TestFilesQueryRepo(t *testing.T) {
	filesQueryRepo := FilesQueryRepo{}
	currentUser, _ := user.Current()
	userHomeDir := "/home/" + currentUser.Username

	t.Run("GetFiles", func(t *testing.T) {
		unixDirPath, _ := valueObject.NewUnixFilePath(userHomeDir)
		_, err := filesQueryRepo.Read(unixDirPath)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
	})

	t.Run("GetFilesWithInvalidDirectory", func(t *testing.T) {
		invalidUnixPath, _ := valueObject.NewUnixFilePath("/aaa/bbb/ccc")
		_, err := filesQueryRepo.Read(invalidUnixPath)
		if err == nil {
			t.Errorf("ExpectedErrorButGotNil")
		}
	})

	t.Run("GetFilesFollowingSymlink", func(t *testing.T) {
		downloadsDirPath, _ := valueObject.NewUnixFilePath(userHomeDir + "/Downloads")
		tmpSymlinkPath, _ := valueObject.NewUnixFilePath(userHomeDir + "/tmpSymlink")

		err := infraHelper.CreateSymlink(
			downloadsDirPath.String(), tmpSymlinkPath.String(), false,
		)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		files, err := filesQueryRepo.Read(tmpSymlinkPath)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		if len(files) == 0 {
			t.Errorf("ExpectedNonEmptyFilesButGotEmpty")
		}

		_ = infraHelper.RemoveSymlink(tmpSymlinkPath.String())
	})

	t.Run("GetSingleFile", func(t *testing.T) {
		unixFilePath, _ := valueObject.NewUnixFilePath(userHomeDir + "/.bashrc")
		_, err := filesQueryRepo.ReadUnique(unixFilePath)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
	})
}
