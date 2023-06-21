package repository

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type AccCmdRepo interface {
	Add(addUser dto.AddUser) error
	Delete(userId valueObject.UserId) error
}
