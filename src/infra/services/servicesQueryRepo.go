package servicesInfra

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"

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

func (repo *ServicesQueryRepo) GetMultiServiceName(
	serviceName valueObject.ServiceName,
	startupFile *valueObject.UnixFilePath,
) (valueObject.ServiceName, error) {
	var startupFilePathStr string

	switch serviceName.String() {
	case "node":
		startupFilePathStr = infraEnvs.PrimaryPublicDir + "/index.js"
	default:
		return "", errors.New("UnknownInstallableMultiService")
	}

	if startupFile != nil {
		startupFilePathStr = startupFile.String()
	}

	startupFileBytes := []byte(startupFilePathStr)
	startupFileHash := md5.Sum(startupFileBytes)
	startupFileHashStr := hex.EncodeToString(startupFileHash[:])
	startupFileShortHashStr := startupFileHashStr[:12]

	svcNameWithSuffix := serviceName.String() + "-" + startupFileShortHashStr
	return valueObject.NewServiceName(svcNameWithSuffix)
}

func (repo ServicesQueryRepo) Get() ([]entity.Service, error) {
	serviceEntities := []entity.Service{}
	return serviceEntities, nil
}

func (repo ServicesQueryRepo) getPpidEntireProcessFamily(
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

func (repo ServicesQueryRepo) getSupervisordServiceMetrics(
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

func (repo ServicesQueryRepo) GetWithMetrics() ([]dto.ServiceWithMetrics, error) {
	servicesWithMetrics := []dto.ServiceWithMetrics{}
	return servicesWithMetrics, nil
}

func (repo ServicesQueryRepo) GetByName(
	name valueObject.ServiceName,
) (serviceEntity entity.Service, err error) {
	return serviceEntity, err
}

func (repo ServicesQueryRepo) GetInstallables() ([]entity.InstallableService, error) {
	return []entity.InstallableService{}, nil
}
