package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
)

type AuthQueryRepo interface {
	IsLoginValid(login dto.Login) bool
}
