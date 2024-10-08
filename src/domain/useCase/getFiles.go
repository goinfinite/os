package useCase

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func GetFiles(
	filesQueryRepo repository.FilesQueryRepo,
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	return filesQueryRepo.Get(unixFilePath)
}
