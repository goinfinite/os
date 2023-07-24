package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type AccQueryRepo interface {
	Get() ([]entity.AccountDetails, error)
	GetByUsername(
		username valueObject.Username,
	) (entity.AccountDetails, error)
	GetById(
		userId valueObject.UserId,
	) (entity.AccountDetails, error)
}
