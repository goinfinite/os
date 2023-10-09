package infra

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	"golang.org/x/exp/slices"

	"github.com/shirou/gopsutil/process"
)

type ServicesQueryRepo struct {
}

func (repo ServicesQueryRepo) ServiceNameAdapter(serviceName string) string {
	switch serviceName {
	case "litespeed":
		return "openlitespeed"
	case "mysqld", "mariadbd", "mariadb-server", "percona-server-mysqld":
		return "mysql"
	case "redis-server":
		return "redis"
	default:
		return serviceName
	}
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
	runningServices, err := repo.runningServiceFactory()
	if err != nil {
		return []entity.Service{}, err
	}

	var runningServicesNames []string
	for _, svc := range runningServices {
		runningServicesNames = append(runningServicesNames, svc.Name.String())
	}

	var notRunningServicesNames []string
	for _, svc := range valueObject.SupportedServiceNames {
		if !slices.Contains(runningServicesNames, svc) {
			notRunningServicesNames = append(notRunningServicesNames, svc)
		}
	}

	var remainingServices []entity.Service
	confFilePath := "/speedia/supervisord.conf"
	for _, svc := range notRunningServicesNames {
		cmd := exec.Command(
			"awk",
			fmt.Sprintf("/%s/{found=1} END{if(!found) exit 1}", svc),
			confFilePath,
		)
		err := cmd.Run()

		svcName, _ := valueObject.NewServiceName(svc)
		svcStatus, _ := valueObject.NewServiceStatus("stopped")
		if err != nil {
			svcStatus, _ = valueObject.NewServiceStatus("uninstalled")
		}

		remainingServices = append(
			remainingServices,
			entity.NewService(
				svcName,
				svcStatus,
				nil,
				nil,
				nil,
				nil,
			),
		)
	}

	return append(runningServices, remainingServices...), nil
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
