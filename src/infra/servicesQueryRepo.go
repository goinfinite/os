package infra

import (
	"errors"
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/shirou/gopsutil/process"
)

type ServicesQueryRepo struct {
}

func (repo ServicesQueryRepo) ServiceNameAdapter(serviceName string) string {
	switch serviceName {
	case "litespeed", "openlitespeed", "lsphp":
		return "php"
	case "mysqld", "mariadbd", "mariadb-server", "percona-server-mysqld":
		return "mysql"
	case "redis-server":
		return "redis"
	default:
		return serviceName
	}
}

func (repo ServicesQueryRepo) getType(
	name valueObject.ServiceName,
) (valueObject.ServiceType, error) {
	svcTypeStr := "runtime"

	switch name.String() {
	case "mysql", "postgresql", "redis", "mongo":
		svcTypeStr = "database"
	}

	return valueObject.NewServiceType(svcTypeStr)
}

func (repo ServicesQueryRepo) runningServiceFactory() ([]entity.Service, error) {
	pids, err := process.Pids()
	if err != nil {
		return []entity.Service{}, errors.New("PidsUnavailable")
	}

	runningStatus, _ := valueObject.NewServiceStatus("running")

	var services []entity.Service
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		procName, err := p.Name()
		if err != nil {
			continue
		}
		procName = repo.ServiceNameAdapter(procName)
		svcName, err := valueObject.NewServiceName(procName)
		if err != nil {
			continue
		}

		svcType, err := repo.getType(svcName)
		if err != nil {
			continue
		}

		uptime, err := p.CreateTime()
		if err != nil {
			continue
		}
		uptimeSeconds := int64(time.Since(time.Unix(uptime/1000, 0)).Seconds())

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}

		memPercent, err := p.MemoryPercent()
		if err != nil {
			continue
		}

		alreadyExists := false
		for i, svc := range services {
			if svc.Name.String() != svcName.String() {
				continue
			}
			alreadyExists = true
			*services[i].Pids = append(*services[i].Pids, uint32(pid))
			if uptimeSeconds > *svc.UptimeSecs {
				*services[i].UptimeSecs = uptimeSeconds
			}
			*services[i].CpuUsagePercent += cpuPercent
			*services[i].MemUsagePercent += memPercent
			continue
		}

		if alreadyExists {
			continue
		}

		var pidUint []uint32
		pidUint = append(pidUint, uint32(pid))

		services = append(
			services,
			entity.NewService(
				svcName,
				svcType,
				runningStatus,
				&pidUint,
				&uptimeSeconds,
				&cpuPercent,
				&memPercent,
			),
		)
	}

	return services, nil
}

func (repo ServicesQueryRepo) Get() ([]entity.Service, error) {
	servicesList := []entity.Service{}

	runningSvcs, err := repo.runningServiceFactory()
	if err != nil {
		return servicesList, err
	}

	var runningSvcNames []string
	for _, svc := range runningSvcs {
		runningSvcNames = append(runningSvcNames, svc.Name.String())
	}

	supervisorConfPath := "/speedia/supervisord.conf"
	supervisorConfContent, err := infraHelper.GetFileContent(supervisorConfPath)
	if err != nil {
		return servicesList, err
	}

	supervisorSvcNameRegex := `(?m)^\[program:(\w{1,64})\]$`
	installedSvcNames := infraHelper.GetRegexCapturingGroups(
		supervisorConfContent,
		supervisorSvcNameRegex,
	)
	if len(installedSvcNames) == 0 {
		return servicesList, errors.New("NoServicesFound")
	}

	supportedSvcNames := maps.Keys(valueObject.SupportedServiceNamesAndAliases)
	for _, svcName := range supportedSvcNames {
		svcName, err := valueObject.NewServiceName(svcName)
		if err != nil {
			continue
		}

		if slices.Contains(runningSvcNames, svcName.String()) {
			continue
		}

		svcType, err := repo.getType(svcName)
		if err != nil {
			continue
		}

		svcStatus, _ := valueObject.NewServiceStatus("uninstalled")
		if slices.Contains(installedSvcNames, svcName.String()) {
			svcStatus, _ = valueObject.NewServiceStatus("stopped")
		}

		servicesList = append(
			servicesList,
			entity.NewService(
				svcName,
				svcType,
				svcStatus,
				nil,
				nil,
				nil,
				nil,
			),
		)
	}

	return append(servicesList, runningSvcs...), nil
}

func (repo ServicesQueryRepo) GetByName(
	name valueObject.ServiceName,
) (entity.Service, error) {
	services, err := repo.Get()
	if err != nil {
		return entity.Service{}, err
	}

	for _, svc := range services {
		svcName := repo.ServiceNameAdapter(svc.Name.String())
		if svcName == name.String() {
			return svc, nil
		}
	}

	return entity.Service{}, errors.New("ServiceNotFound")
}
