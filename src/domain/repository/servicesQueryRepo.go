package repository

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesQueryRepo interface {
	Get() ([]entity.Service, error)
	GetByName(name valueObject.ServiceName) (entity.Service, error)
	GetInstallables() ([]entity.InstallableService, error)
}
