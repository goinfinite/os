package servicesInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesCmdRepo struct {
}

func (repo ServicesCmdRepo) AddInstallable(
	addDto dto.AddInstallableService,
) error {
	err := AddInstallable(addDto)
	if err != nil {
		return err
	}

	err = SupervisordFacade{}.Reload()
	if err != nil {
		return err
	}

	return nil
}

func (repo ServicesCmdRepo) AddCustom(
	addDto dto.CreateCustomService,
) error {
	err := AddCustom(addDto)
	if err != nil {
		return err
	}

	err = SupervisordFacade{}.Reload()
	if err != nil {
		return err
	}

	return nil
}

func (repo ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	return SupervisordFacade{}.Start(name)
}

func (repo ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	return SupervisordFacade{}.Stop(name)
}

func (repo ServicesCmdRepo) Restart(name valueObject.ServiceName) error {
	return SupervisordFacade{}.Restart(name)
}

func (repo ServicesCmdRepo) Update(
	updateDto dto.UpdateService,
) error {
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

	err = SupervisordFacade{}.Reload()
	if err != nil {
		return err
	}

	return nil
}

func (repo ServicesCmdRepo) Uninstall(
	name valueObject.ServiceName,
) error {
	err := repo.Stop(name)
	if err != nil {
		return err
	}

	err = Uninstall(name)
	if err != nil {
		return err
	}

	err = SupervisordFacade{}.Reload()
	if err != nil {
		return err
	}

	return nil
}
