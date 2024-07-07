package servicesInfra

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
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
			log.Printf("[%s] %s", serviceModel.Name, err.Error())
			continue
		}

		servicesEntities = append(servicesEntities, serviceEntity)
	}

	return servicesEntities, nil
}

func (repo *ServicesQueryRepo) ReadByName(
	name valueObject.ServiceName,
) (serviceEntity entity.InstalledService, err error) {
	var serviceModel dbModel.InstalledService
	queryResult := repo.persistentDbSvc.Handler.
		Where("name = ?", name.String()).
		Limit(1).
		Find(&serviceModel)
	if queryResult.Error != nil {
		return serviceEntity, err
	}

	if queryResult.RowsAffected == 0 {
		return serviceEntity, errors.New("ServiceNotFound")
	}

	serviceEntity, err = serviceModel.ToEntity()
	if err != nil {
		return serviceEntity, err
	}

	return serviceEntity, nil
}

func (repo *ServicesQueryRepo) getPidProcessFamily(pid int32) ([]*process.Process, error) {
	processFamily := []*process.Process{}

	pidProcess, err := process.NewProcess(pid)
	if err != nil {
		return processFamily, err
	}

	processFamily = append(processFamily, pidProcess)

	childrenPidProcesses, err := pidProcess.Children()
	if err != nil || len(childrenPidProcesses) == 0 {
		return processFamily, nil
	}

	for _, childPidProcess := range childrenPidProcesses {
		grandChildrenPidProcesses, err := repo.getPidProcessFamily(
			childPidProcess.Pid,
		)
		if err != nil || len(grandChildrenPidProcesses) == 0 {
			continue
		}

		processFamily = append(processFamily, grandChildrenPidProcesses...)
	}

	return processFamily, nil
}

func (repo *ServicesQueryRepo) getPidMetrics(
	mainPid int32,
) (serviceMetrics valueObject.ServiceMetrics, err error) {
	pidProcesses, err := repo.getPidProcessFamily(mainPid)
	if err != nil {
		return serviceMetrics, err
	}

	if len(pidProcesses) == 0 {
		return serviceMetrics, nil
	}

	uptimeMilliseconds, err := pidProcesses[0].CreateTime()
	if err != nil {
		return serviceMetrics, err
	}
	nowMilliseconds := time.Now().UTC().UnixMilli()
	uptimeSecs := (nowMilliseconds - uptimeMilliseconds) / 1000

	cpuPercent := float64(0.0)
	memPercent := float32(0.0)

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

	cpuPercent = math.Round(cpuPercent*100) / 100
	memPercent = float32(math.Round(float64(memPercent)*100) / 100)

	serviceMetrics = valueObject.NewServiceMetrics(
		pids, uptimeSecs, cpuPercent, memPercent,
	)

	return serviceMetrics, nil
}

func (repo *ServicesQueryRepo) ReadWithMetrics() ([]dto.InstalledServiceWithMetrics, error) {
	servicesWithMetrics := []dto.InstalledServiceWithMetrics{}

	servicesEntities, err := repo.Read()
	if err != nil {
		return servicesWithMetrics, err
	}
	serviceNameServiceEntityMap := map[string]entity.InstalledService{}
	for _, serviceEntity := range servicesEntities {
		serviceNameServiceEntityMap[serviceEntity.Name.String()] = serviceEntity
	}

	supervisorStatus, _ := infraHelper.RunCmd("supervisorctl", "status")
	if len(supervisorStatus) == 0 {
		return servicesWithMetrics, errors.New("GetSupervisorStatusError")
	}

	// # supervisorctl status
	// cron                             RUNNING   pid 2, uptime 0:00:11
	// nginx                            RUNNING   pid 24, uptime 0:00:10
	// os-api                           RUNNING   pid 3, uptime 0:00:11
	supervisorStatusLines := strings.Split(supervisorStatus, "\n")
	if len(supervisorStatusLines) == 0 {
		return servicesWithMetrics, errors.New("SupervisorStatusEmpty")
	}

	for _, supervisorStatusLine := range supervisorStatusLines {
		if supervisorStatusLine == "" {
			continue
		}

		supervisorStatusLineParts := strings.Fields(supervisorStatusLine)
		if len(supervisorStatusLineParts) != 6 {
			continue
		}

		rawServiceName := supervisorStatusLineParts[0]
		serviceName, err := valueObject.NewServiceName(rawServiceName)
		if err != nil {
			continue
		}

		serviceEntity, exists := serviceNameServiceEntityMap[serviceName.String()]
		if !exists {
			continue
		}

		rawServiceStatus := supervisorStatusLineParts[1]
		serviceStatus, err := valueObject.NewServiceStatus(rawServiceStatus)
		if err != nil {
			continue
		}

		if serviceStatus.String() != "running" {
			serviceWithMetrics := dto.NewInstalledServiceWithMetrics(serviceEntity, nil)
			servicesWithMetrics = append(servicesWithMetrics, serviceWithMetrics)
			continue
		}

		rawServicePid := supervisorStatusLineParts[3]
		rawServicePid = strings.Trim(rawServicePid, ",")
		servicePidInt, err := strconv.ParseInt(rawServicePid, 10, 32)
		if err != nil {
			continue
		}

		serviceMetrics, err := repo.getPidMetrics(int32(servicePidInt))
		if err != nil {
			continue
		}

		serviceWithMetrics := dto.NewInstalledServiceWithMetrics(
			serviceEntity, &serviceMetrics,
		)

		servicesWithMetrics = append(servicesWithMetrics, serviceWithMetrics)
	}

	return servicesWithMetrics, nil
}

func (repo *ServicesQueryRepo) ReadInstallables() ([]entity.InstallableService, error) {
	return []entity.InstallableService{}, nil
}
