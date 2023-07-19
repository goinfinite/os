package repository

import (
	"github.com/speedianet/sam/src/domain/valueObject"
)

type DatabaseCmdRepo interface {
	Add(dbName valueObject.DatabaseName) error
}
