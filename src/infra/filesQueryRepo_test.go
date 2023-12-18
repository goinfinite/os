package infra

import (
	"testing"

	"github.com/speedianet/os/src/domain/valueObject"
)

func TestFilesQueryRepo(t *testing.T) {
	unixFilePath, _ := valueObject.NewUnixFilePath("/etc/timezone")
	invalidUnixPath, _ := valueObject.NewUnixFilePath("/aaa/bbb/ccc")

	filesQueryRepo := FilesQueryRepo{}

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
		unixFilePath, _ := valueObject.NewUnixFilePath("/etc/ssl")
		_, err := filesQueryRepo.Get(unixFilePath)
		if err != nil {
			t.Errorf("FilesQueryRepo.Get() should not return: %s", err)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		_, err := filesQueryRepo.Exists(unixFilePath)
		if err != nil {
			t.Errorf("FilesQueryRepo.Exists() should not return: %s", err)
		}
	})
}
