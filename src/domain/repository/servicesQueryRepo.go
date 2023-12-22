package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesQueryRepo interface {
	Get() ([]entity.Service, error)
	GetWithMetrics() ([]dto.ServiceWithMetrics, error)
	GetByName(name valueObject.ServiceName) (entity.Service, error)
	GetInstallables() ([]entity.InstallableService, error)
}
