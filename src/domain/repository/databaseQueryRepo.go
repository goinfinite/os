package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DatabaseQueryRepo interface {
	Read() ([]entity.Database, error)
	ReadByName(dbName valueObject.DatabaseName) (entity.Database, error)
}
