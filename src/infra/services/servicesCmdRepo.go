package servicesInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesCmdRepo struct {
	supervisordFacade SupervisordFacade
}

func NewServicesCmdRepo() *ServicesCmdRepo {
	return &ServicesCmdRepo{supervisordFacade: SupervisordFacade{}}
}

func (repo *ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	return repo.supervisordFacade.Start(name)
}

func (repo *ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	return repo.supervisordFacade.Stop(name)
}

func (repo *ServicesCmdRepo) Restart(name valueObject.ServiceName) error {
	return repo.supervisordFacade.Restart(name)
}

func (repo *ServicesCmdRepo) Reload() error {
	return repo.supervisordFacade.Reload()
}

func (repo *ServicesCmdRepo) CreateInstallable(
	createDto dto.CreateInstallableService,
) error {
	err := CreateInstallable(createDto)
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) CreateCustom(createDto dto.CreateCustomService) error {
	err := CreateCustom(createDto)
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) Update(updateDto dto.UpdateService) error {
	err := repo.Stop(updateDto.Name)
	if err != nil {
		return err
	}

	serviceEntity, err := ServicesQueryRepo{}.GetByName(updateDto.Name)
	if err != nil {
		return err
	}

	err = Update(serviceEntity, updateDto)
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) Uninstall(name valueObject.ServiceName) error {
	err := repo.Stop(name)
	if err != nil {
		return err
	}

	err = Uninstall(name)
	if err != nil {
		return err
	}

	return repo.Reload()
}
