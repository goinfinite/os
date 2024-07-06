package servicesInfra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type ServicesCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServicesCmdRepo {
	return &ServicesCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd("supervisorctl", "start", name.String())
	return err
}

func (repo *ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd("supervisorctl", "stop", name.String())
	return err
}

func (repo *ServicesCmdRepo) Restart(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd("supervisorctl", "restart", name.String())
	return err
}

func (repo *ServicesCmdRepo) Reload() error {
	_, err := infraHelper.RunCmd("supervisorctl", "reload")
	return err
}

func (repo *ServicesCmdRepo) CreateInstallable(
	createDto dto.CreateInstallableService,
) error {
	return repo.Reload()
}

func (repo *ServicesCmdRepo) CreateCustom(createDto dto.CreateCustomService) error {
	customNature, _ := valueObject.NewServiceNature("custom")

	installedServiceModel := dbModel.NewInstalledService(
		createDto.Name.String(),
		customNature.String(),
		createDto.Type.String(),
		createDto.Version.String(),
		createDto.Command.String(),
		createDto.Envs,
		createDto.PortBindings,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	err := repo.persistentDbSvc.Handler.Create(&installedServiceModel).Error
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) Update(updateDto dto.UpdateService) error {
	return repo.Reload()
}

func (repo *ServicesCmdRepo) Delete(name valueObject.ServiceName) error {
	return repo.Reload()
}
