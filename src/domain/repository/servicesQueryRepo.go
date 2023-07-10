package repository

import (
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
)

type ServicesQueryRepo interface {
	Get() ([]entity.Service, error)
	GetByName(name valueObject.ServiceName) (entity.Service, error)
}
