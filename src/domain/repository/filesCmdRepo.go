package repository

import "github.com/speedianet/os/src/domain/dto"

type FilesCmdRepo interface {
	Add(dto.AddUnixFile) error
}
