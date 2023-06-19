package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type AccQueryRepo interface {
	GetAccountDetailsByUsername(
		username valueObject.Username,
	) entity.AccountDetails
}
