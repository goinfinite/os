package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DatabaseQueryRepo interface {
	Get() ([]entity.Database, error)
	GetByName(dbName valueObject.DatabaseName) (entity.Database, error)
}
