package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
)

type DatabaseQueryRepo interface {
	Get() ([]entity.Database, error)
}
