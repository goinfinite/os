package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type DatabaseQueryRepo interface {
	Get() ([]entity.Database, error)
	GetByName(dbName valueObject.DatabaseName) (entity.Database, error)
}
