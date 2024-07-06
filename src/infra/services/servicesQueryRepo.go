package servicesInfra

import (
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"

	"github.com/shirou/gopsutil/process"
)

type ServicesQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServicesQueryRepo {
	return &ServicesQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ServicesQueryRepo) Read() ([]entity.InstalledService, error) {
	servicesEntities := []entity.InstalledService{}

	servicesModels := []dbModel.InstalledService{}
	err := repo.persistentDbSvc.Handler.
		Find(&servicesModels).Error
	if err != nil {
		return servicesEntities, err
	}

	for _, serviceModel := range servicesModels {
		serviceEntity, err := serviceModel.ToEntity()
		if err != nil {
			log.Printf("InstalledServiceModelToEntityError: %s", err.Error())
			continue
		}

		servicesEntities = append(servicesEntities, serviceEntity)
	}

	return servicesEntities, nil
}

func (repo *ServicesQueryRepo) ReadByName(
	name valueObject.ServiceName,
) (serviceEntity entity.InstalledService, err error) {
	serviceModel := dbModel.InstalledService{}
	err = repo.persistentDbSvc.Handler.
		Where("name = ?", name.String()).
		First(&serviceModel).Error
	if err != nil {
		return serviceEntity, err
	}

	serviceEntity, err = serviceModel.ToEntity()
	if err != nil {
		return serviceEntity, err
	}

	return serviceEntity, nil
}

func (repo *ServicesQueryRepo) getPpidEntireProcessFamily(
	ppid int32,
) ([]*process.Process, error) {
	ppidProcesses := []*process.Process{}

	ppidProcess, err := process.NewProcess(ppid)
	if err != nil {
		return ppidProcesses, err
	}

	ppidProcesses = append(ppidProcesses, ppidProcess)

	childrenPidProcesses, _ := ppidProcess.Children()
	if len(childrenPidProcesses) == 0 {
		return ppidProcesses, nil
	}

	for _, childPidProcess := range childrenPidProcesses {
		grandChildrenPidProcesses, _ := repo.getPpidEntireProcessFamily(
			childPidProcess.Pid,
		)
		if len(grandChildrenPidProcesses) == 0 {
			continue
		}

		ppidProcesses = append(ppidProcesses, grandChildrenPidProcesses...)
	}

	return ppidProcesses, nil
}

func (repo *ServicesQueryRepo) getSupervisordServiceMetrics(
	mainPid int32,
	uptimeSecs int64,
) (valueObject.ServiceMetrics, error) {
	supervisordServiceMetrics := valueObject.ServiceMetrics{}

	cpuPercent := float64(0.0)
	memPercent := float32(0.0)

	pidProcesses, err := repo.getPpidEntireProcessFamily(mainPid)
	if err != nil {
		return supervisordServiceMetrics, err
	}

	pids := []uint32{}
	for _, process := range pidProcesses {
		pidCpuPercent, err := process.CPUPercent()
		if err != nil {
			continue
		}

		pidMemPercent, err := process.MemoryPercent()
		if err != nil {
			continue
		}

		cpuPercent += pidCpuPercent
		memPercent += pidMemPercent

		pids = append(pids, uint32(process.Pid))
	}

	serviceMetrics := valueObject.NewServiceMetrics(
		pids,
		uptimeSecs,
		cpuPercent,
		memPercent,
	)

	return serviceMetrics, nil
}

func (repo *ServicesQueryRepo) ReadWithMetrics() ([]dto.InstalledServiceWithMetrics, error) {
	servicesWithMetrics := []dto.InstalledServiceWithMetrics{}
	return servicesWithMetrics, nil
}

func (repo *ServicesQueryRepo) ReadInstallables() ([]entity.InstallableService, error) {
	return []entity.InstallableService{}, nil
}
