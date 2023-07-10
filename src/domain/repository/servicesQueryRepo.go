package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
)

type ServicesQueryRepo interface {
	Get() ([]entity.Service, error)
}
