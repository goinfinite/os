package infra

import (
	"testing"

	"github.com/speedianet/os/src/domain/valueObject"
)

func TestFilesQueryRepo(t *testing.T) {
	unixFilePath, _ := valueObject.NewUnixFilePath("/etc/timezone")

	filesQueryRepo := FilesQueryRepo{}

	t.Run("Get", func(t *testing.T) {
		_, err := filesQueryRepo.Get(unixFilePath)
		if err != nil {
			t.Errorf("FilesQueryRepo.Get() should not return: %s", err)
		}
	})

	t.Run("Get (many)", func(t *testing.T) {
		unixFilePath, _ := valueObject.NewUnixFilePath("/home/mmp/.vscode")
		_, err := filesQueryRepo.Get(unixFilePath)
		if err != nil {
			t.Errorf("FilesQueryRepo.Get() should not return: %s", err)
		}
	})
}
