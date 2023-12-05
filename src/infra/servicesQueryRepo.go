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

	envsRegex := regexp.MustCompile(
		`(?P<envName>[a-zA-Z0-9_\-]{1,256})="(?P<envValue>[^"]{1,256})"`,
	)
	envsMatches := envsRegex.FindAllStringSubmatch(envs, -1)

	for _, match := range envsMatches {
		if len(match) != 3 {
			continue
		}

		key := strings.TrimSpace(match[1])
		value := strings.TrimSpace(match[2])
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
			`(?m)^\[program:(` + svcNameRegex + `)\]`,
		)
		svcNameMatches := svcNameRegex.FindStringSubmatch(svcConfigBlock)
		if len(svcNameMatches) == 0 {
			continue
		}

		svcName, err := valueObject.NewServiceName(svcNameMatches[1])
		if err != nil {
			continue
		}

		svcCmdRegex := regexp.MustCompile(`(?m)^command\=(.*)$`)
		svcCmdMatches := svcCmdRegex.FindStringSubmatch(svcConfigBlock)
		if len(svcCmdMatches) == 0 {
			continue
		}

		svcCmd, err := valueObject.NewUnixCommand(svcCmdMatches[1])
		if err != nil {
			continue
		}

		svcEnvsRegex := regexp.MustCompile(`(?m)^environment\=(.*)$`)
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

		svcPortsStr, exists := svcEnvs["SVC_PORTS"]
		if !exists {
			continue
		}

		svcPortsParts := strings.Split(svcPortsStr, ",")
		if len(svcPortsParts) == 0 {
			continue
		}
		svcPorts := []valueObject.NetworkPort{}
		for _, svcPortStr := range svcPortsParts {
			svcPort, err := valueObject.NewNetworkPort(svcPortStr)
			if err != nil {
				continue
			}
			svcPorts = append(svcPorts, svcPort)
		}

		servicesList = append(
			servicesList,
			entity.NewService(
				svcName,
				svcType,
				svcStatus,
				&svcCmd,
				svcPorts,
				[]uint32{},
				nil,
				nil,
				nil,
			),
		)
	}

	if len(servicesList) == 0 {
		return servicesList, errors.New("GetInstalledServicesFailed")
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
			[]valueObject.NetworkPort{},
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
	servicesList []entity.Service,
) ([]entity.Service, error) {
	pids, err := process.Pids()
	if err != nil {
		return servicesList, errors.New("ServicePidsUnavailable")
	}

	svcRunningStatus, _ := valueObject.NewServiceStatus("running")

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
		serviceEntityIndex := 0
		for svcIndex, serviceFromList := range servicesList {
			if serviceFromList.Name.String() != svcName.String() {
				continue
			}

			serviceEntity = serviceFromList
			serviceEntityIndex = svcIndex
		}

		if serviceEntity.Name.String() == "" {
			continue
		}

		serviceEntity.Status = svcRunningStatus

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

		servicesList[serviceEntityIndex] = serviceEntity
	}

	return servicesList, nil
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

func (repo ServicesQueryRepo) GetInstallables() ([]entity.InstallableService, error) {
	phpService := entity.NewInstallableService(
		valueObject.NewServiceNamePanic("php"),
		valueObject.NewServiceTypePanic("runtime"),
		[]valueObject.ServiceVersion{
			valueObject.NewServiceVersionPanic("8.2"),
			valueObject.NewServiceVersionPanic("8.1"),
			valueObject.NewServiceVersionPanic("8.0"),
			valueObject.NewServiceVersionPanic("7.4"),
			valueObject.NewServiceVersionPanic("7.3"),
			valueObject.NewServiceVersionPanic("7.2"),
			valueObject.NewServiceVersionPanic("7.1"),
			valueObject.NewServiceVersionPanic("7.0"),
			valueObject.NewServiceVersionPanic("5.6"),
		},
	)

	nodeService := entity.NewInstallableService(
		valueObject.NewServiceNamePanic("node"),
		valueObject.NewServiceTypePanic("runtime"),
		[]valueObject.ServiceVersion{
			valueObject.NewServiceVersionPanic("21"),
			valueObject.NewServiceVersionPanic("20"),
			valueObject.NewServiceVersionPanic("19"),
			valueObject.NewServiceVersionPanic("18"),
			valueObject.NewServiceVersionPanic("17"),
			valueObject.NewServiceVersionPanic("16"),
			valueObject.NewServiceVersionPanic("15"),
			valueObject.NewServiceVersionPanic("14"),
			valueObject.NewServiceVersionPanic("13"),
			valueObject.NewServiceVersionPanic("12"),
			valueObject.NewServiceVersionPanic("11"),
			valueObject.NewServiceVersionPanic("10"),
			valueObject.NewServiceVersionPanic("9"),
			valueObject.NewServiceVersionPanic("8"),
			valueObject.NewServiceVersionPanic("7"),
			valueObject.NewServiceVersionPanic("6"),
			valueObject.NewServiceVersionPanic("5"),
			valueObject.NewServiceVersionPanic("4"),
		},
	)

	mariadbService := entity.NewInstallableService(
		valueObject.NewServiceNamePanic("mariadb"),
		valueObject.NewServiceTypePanic("database"),
		[]valueObject.ServiceVersion{
			valueObject.NewServiceVersionPanic("10.11"),
			valueObject.NewServiceVersionPanic("10.6"),
		},
	)

	redisService := entity.NewInstallableService(
		valueObject.NewServiceNamePanic("redis"),
		valueObject.NewServiceTypePanic("database"),
		[]valueObject.ServiceVersion{
			valueObject.NewServiceVersionPanic("7.2"),
			valueObject.NewServiceVersionPanic("7.0"),
			valueObject.NewServiceVersionPanic("6.2"),
		},
	)

	return []entity.InstallableService{
		phpService,
		nodeService,
		mariadbService,
		redisService,
	}, nil
}
