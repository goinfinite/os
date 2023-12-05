package infra

import (
	"mime"
	"strings"
	"testing"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestFilesCmdRepo(t *testing.T) {
	filesCmdRepo := FilesCmdRepo{}

	t.Run("AddValidUnixFile", func(t *testing.T) {
		unixFilePath := valueObject.NewUnixFilePathPanic("/home/mmp/filesCmdRepoTest.txt")

		unixFileExtension, _ := unixFilePath.GetFileExtension()
		mimeTypeWithCharset := mime.TypeByExtension("." + unixFileExtension.String())
		mimeTypeOnly := strings.Split(mimeTypeWithCharset, ";")[0]

		unixFileName, _ := unixFilePath.GetFileName()

		addUnixFile := dto.NewAddUnixFile(
			valueObject.NewMimeTypePanic(mimeTypeOnly),
			unixFileName,
			unixFilePath,
			valueObject.NewUnixFilePermissionsPanic("0777"),
		)

		err := filesCmdRepo.Add(addUnixFile)
		if err != nil {
			t.Errorf("UnexpectedError: %v", err)
		}
	})
}
