package filesInfra

import (
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

	t.Run("ReadFollowingSymlink", func(t *testing.T) {
		downloadsDirPath, _ := tkValueObject.NewUnixAbsoluteFilePath(userHomeDir+"/Downloads", false)
		tmpSymlinkPath, _ := tkValueObject.NewUnixAbsoluteFilePath(userHomeDir+"/tmpSymlink", false)
		requestDto := dto.ReadFilesRequest{
			SourcePath: tmpSymlinkPath,
		}

		err := fileClerk.CreateSymlink(
			downloadsDirPath.String(), tmpSymlinkPath.String(), false,
		)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}

		responseDto, err := filesQueryRepo.Read(requestDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
		if len(responseDto.Files) == 0 {
			t.Errorf("ExpectedNonEmptyFilesButGotEmpty")
		}

		_ = fileClerk.RemoveSymlink(tmpSymlinkPath.String())
	})

	t.Run("ReadFirstFile", func(t *testing.T) {
		unixFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(userHomeDir+"/.bashrc", false)
		_, err := filesQueryRepo.ReadFirst(unixFilePath)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %s", err.Error())
		}
	})
}
