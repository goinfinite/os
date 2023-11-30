package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func GetFiles(
	filesQueryRepo repository.FilesQueryRepo,
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	return filesQueryRepo.Get(unixFilePath)
}
