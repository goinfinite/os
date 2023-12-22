package infra

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/speedianet/os/src/domain/valueObject"
)

func TestFilesQueryRepo(t *testing.T) {
	filesQueryRepo := FilesQueryRepo{}

	currentUser, _ := user.Current()
	fileBasePathStr := fmt.Sprintf("/home/%s", currentUser.Username)

	unixFilePath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/.gitconfig")
	invalidUnixPath, _ := valueObject.NewUnixFilePath("/aaa/bbb/ccc")
	unixDirPath, _ := valueObject.NewUnixFilePath(fileBasePathStr + "/Downloads")

	t.Run("Get", func(t *testing.T) {
		_, err := filesQueryRepo.Get(unixFilePath)
		if err != nil {
			t.Errorf("FilesQueryRepo.Get() should not return: %s", err)
		}
	})

	t.Run("InvalidGet", func(t *testing.T) {
		_, err := filesQueryRepo.Get(invalidUnixPath)
		if err == nil {
			t.Errorf("FilesQueryRepo.Get() should throw an error")
		}
	})

	t.Run("Get (many)", func(t *testing.T) {
		_, err := filesQueryRepo.Get(unixDirPath)
		if err != nil {
			t.Errorf("FilesQueryRepo.Get() should not return: %s", err)
		}
	})

	t.Run("GetOnly", func(t *testing.T) {
		_, err := filesQueryRepo.GetOnly(unixDirPath)
		if err != nil {
			t.Errorf("FilesQueryRepo.GetOnly() should not return: %s", err)
		}
	})
}
