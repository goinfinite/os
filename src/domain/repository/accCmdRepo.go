package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
)

type AccCmdRepo interface {
	Add(addUser dto.AddUser) error
}
