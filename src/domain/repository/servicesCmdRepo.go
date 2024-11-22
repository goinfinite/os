package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ServicesCmdRepo interface {
	CreateInstallable(dto.CreateInstallableService) (valueObject.ServiceName, error)
	CreateCustom(dto.CreateCustomService) error
	Update(dto.UpdateService) error
	Delete(valueObject.ServiceName) error
	RefreshInstallableItems() error
}
