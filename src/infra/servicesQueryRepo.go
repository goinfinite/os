package infra

import (
	"errors"
	"regexp"
	"strings"
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

func (repo ServicesQueryRepo) parseServiceEnvs(envs string) map[string]string {
	envsMap := map[string]string{}

	envsStr := strings.Split(envs, ",")
	for _, envStr := range envsStr {
		env := strings.Split(envStr, "=")
		if len(env) != 2 {
			continue
		}

		key := strings.TrimSpace(env[0])
		value := strings.TrimSpace(env[1])
		envsMap[key] = value
	}

	return envsMap
}

func (repo ServicesQueryRepo) getInstalledServices() ([]entity.Service, error) {
	servicesList := []entity.Service{}

	supervisorConfPath := "/speedia/supervisord.conf"
	supervisorConfContent, err := infraHelper.GetFileContent(supervisorConfPath)
	if err != nil {
		return servicesList, err
	}

	svcNameRegex := strings.TrimLeft(valueObject.ServiceNameRegex, "^")
	svcNameRegex = strings.TrimRight(svcNameRegex, "$")

	svcConfigBlocksRegex := regexp.MustCompile(
		`(?m)^\[program:` + svcNameRegex + `\]\n(?:[^\[]+\n)*`,
	)
	svcConfigBlocks := svcConfigBlocksRegex.FindAllString(supervisorConfContent, -1)
	if len(svcConfigBlocks) == 0 {
		return servicesList, errors.New("NoServicesFound")
	}

	svcStatus, _ := valueObject.NewServiceStatus("stopped")

	for _, svcConfigBlock := range svcConfigBlocks {
		svcNameRegex := regexp.MustCompile(
			`^\[program:(` + svcNameRegex + `)\]`,
		)
		svcNameMatches := svcNameRegex.FindStringSubmatch(svcConfigBlock)
		if len(svcNameMatches) == 0 {
			continue
		}

		svcName, err := valueObject.NewServiceName(svcNameMatches[1])
		if err != nil {
			continue
		}

		svcCmdRegex := regexp.MustCompile(
			`^command=(.*)$`,
		)
		svcCmdMatches := svcCmdRegex.FindStringSubmatch(svcConfigBlock)
		if len(svcCmdMatches) == 0 {
			continue
		}

		svcCmd, err := valueObject.NewUnixCommand(svcCmdMatches[1])
		if err != nil {
			continue
		}

		svcEnvsRegex := regexp.MustCompile(
			`^environment=(.*)$`,
		)
		svcEnvsMatches := svcEnvsRegex.FindStringSubmatch(svcConfigBlock)
		if len(svcEnvsMatches) == 0 {
			continue
		}

		svcEnvs := repo.parseServiceEnvs(svcEnvsMatches[1])

		svcTypeStr, exists := svcEnvs["SVC_TYPE"]
		if !exists {
			continue
		}

		svcType, err := valueObject.NewServiceType(svcTypeStr)
		if err != nil {
			continue
		}

		svcPortStr, exists := svcEnvs["SVC_PORT"]
		if !exists {
			continue
		}

		svcPort, err := valueObject.NewNetworkPort(svcPortStr)
		if err != nil {
			continue
		}

		servicesList = append(
			servicesList,
			entity.NewService(
				svcName,
				svcType,
				svcStatus,
				&svcCmd,
				&svcPort,
				[]uint32{},
				nil,
				nil,
				nil,
			),
		)
	}

	return servicesList, nil
}

func (repo ServicesQueryRepo) getNativeServices() ([]entity.Service, error) {
	servicesList := []entity.Service{}

	svcStatus, _ := valueObject.NewServiceStatus("uninstalled")

	nativeSvcNames := maps.Keys(valueObject.NativeSvcNamesWithAliases)
	for _, nativeSvcName := range nativeSvcNames {
		svcName, err := valueObject.NewServiceName(nativeSvcName)
		if err != nil {
			continue
		}

		svcType, _ := valueObject.NewServiceType("runtime")
		switch svcName.String() {
		case "mysql", "redis", "postgres", "mongo", "memcached", "elasticsearch":
			svcType, _ = valueObject.NewServiceType("database")
		}

		svcEntity := entity.NewService(
			svcName,
			svcType,
			svcStatus,
			nil,
			nil,
			[]uint32{},
			nil,
			nil,
			nil,
		)

		servicesList = append(servicesList, svcEntity)
	}

	return servicesList, nil
}

func (repo ServicesQueryRepo) addServicesMetrics(
	installedServices []entity.Service,
) ([]entity.Service, error) {
	pids, err := process.Pids()
	if err != nil {
		return installedServices, errors.New("ServicePidsUnavailable")
	}

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		procName, err := p.Name()
		if err != nil {
			continue
		}

		svcName, err := valueObject.NewServiceName(procName)
		if err != nil {
			continue
		}

		serviceEntity := entity.Service{}
		for _, installedSvc := range installedServices {
			if installedSvc.Name.String() != svcName.String() {
				continue
			}

			serviceEntity = installedSvc
		}

		var pidUint []uint32
		pidUint = append(pidUint, uint32(pid))
		serviceEntity.Pids = pidUint

		uptime, err := p.CreateTime()
		if err != nil {
			continue
		}
		uptimeSeconds := int64(time.Since(time.Unix(uptime/1000, 0)).Seconds())
		serviceEntity.UptimeSecs = &uptimeSeconds

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}
		serviceEntity.CpuUsagePercent = &cpuPercent

		memPercent, err := p.MemoryPercent()
		if err != nil {
			continue
		}
		serviceEntity.MemUsagePercent = &memPercent
	}

	return installedServices, nil
}

func (repo ServicesQueryRepo) Get() ([]entity.Service, error) {
	servicesList := []entity.Service{}

	installedSvcs, err := repo.getInstalledServices()
	if err != nil {
		return servicesList, err
	}
	servicesList = append(servicesList, installedSvcs...)

	installedSvcsNames := []string{}
	for _, installedSvc := range installedSvcs {
		installedSvcsNames = append(installedSvcsNames, installedSvc.Name.String())
	}

	nativeSvcs, err := repo.getNativeServices()
	if err != nil {
		return servicesList, err
	}

	for _, nativeSvc := range nativeSvcs {
		if slices.Contains(installedSvcsNames, nativeSvc.Name.String()) {
			continue
		}

		servicesList = append(servicesList, nativeSvc)
	}

	svcsWithMetrics, err := repo.addServicesMetrics(servicesList)
	if err != nil {
		return servicesList, err
	}

	return svcsWithMetrics, nil
}

func (repo ServicesQueryRepo) GetByName(
	name valueObject.ServiceName,
) (entity.Service, error) {
	service := entity.Service{}

	services, err := repo.Get()
	if err != nil {
		return service, err
	}

	for _, svc := range services {
		if svc.Name.String() == name.String() {
			return svc, nil
		}
	}

	return service, errors.New("ServiceNotFound")
}
