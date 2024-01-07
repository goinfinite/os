package servicesInfra

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"

	"github.com/shirou/gopsutil/process"
)

type ServicesQueryRepo struct {
}

type supervisordService struct {
	Name            string
	Status          string
	MainPid         int32
	UptimeInSeconds int64
}

func (repo ServicesQueryRepo) getSupervisordServices() ([]supervisordService, error) {
	supervisordServices := []supervisordService{}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	httpResponse, err := httpClient.Get("http://127.0.0.1:9001/program/list")
	if err != nil {
		return supervisordServices, errors.New("SupervisordApiUnavailable")
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return supervisordServices, errors.New("SupervisordApiResponseError")
	}

	var parsedResponse []interface{}
	err = json.Unmarshal(responseBody, &parsedResponse)
	if err != nil {
		return supervisordServices, errors.New("SupervisordApiDecodeResponseError")
	}

	nowEpoch := time.Now().Unix()

	for _, svcDetails := range parsedResponse {
		svcDetailsMap, assertOk := svcDetails.(map[string]interface{})
		if !assertOk {
			continue
		}

		svcName, assertOk := svcDetailsMap["name"].(string)
		if !assertOk {
			continue
		}

		svcStatus, assertOk := svcDetailsMap["statename"].(string)
		if !assertOk {
			continue
		}
		svcStatus = strings.ToLower(svcStatus)

		svcStart, assertOk := svcDetailsMap["start"].(float64)
		if !assertOk {
			continue
		}
		svcUptime := nowEpoch - int64(svcStart)

		svcPid, assertOk := svcDetailsMap["pid"].(float64)
		if !assertOk {
			continue
		}

		supervisordServices = append(
			supervisordServices,
			supervisordService{
				Name:            svcName,
				Status:          svcStatus,
				MainPid:         int32(svcPid),
				UptimeInSeconds: svcUptime,
			},
		)
	}

	return supervisordServices, nil
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

func (repo ServicesQueryRepo) Get() ([]entity.Service, error) {
	servicesList := []entity.Service{}

	supervisordConfPath := "/speedia/supervisord.conf"
	supervisordConfContent, err := infraHelper.GetFileContent(supervisordConfPath)
	if err != nil {
		return servicesList, err
	}

	svcNameRegex := strings.TrimLeft(valueObject.ServiceNameRegex, "^")
	svcNameRegex = strings.TrimRight(svcNameRegex, "$")

	svcConfigBlocksRegex := regexp.MustCompile(
		`(?m)^\[program:` + svcNameRegex + `\]\n(?:[^\[]+\n)*`,
	)
	svcConfigBlocks := svcConfigBlocksRegex.FindAllString(supervisordConfContent, -1)
	if len(svcConfigBlocks) == 0 {
		return servicesList, errors.New("NoServicesFound")
	}

	supervisordServicesFromApi, err := repo.getSupervisordServices()
	if err != nil {
		return servicesList, err
	}
	runningServicesNames := []string{}
	for _, supervisordService := range supervisordServicesFromApi {
		if supervisordService.Status != "running" {
			continue
		}

		runningServicesNames = append(runningServicesNames, supervisordService.Name)
	}
	stoppedStatus, _ := valueObject.NewServiceStatus("stopped")
	runningStatus, _ := valueObject.NewServiceStatus("running")

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

		svcNatureStr, exists := svcEnvs["SVC_NATURE"]
		if !exists {
			continue
		}
		svcNature, err := valueObject.NewServiceNature(svcNatureStr)
		if err != nil {
			continue
		}

		svcVersionStr, exists := svcEnvs["SVC_VERSION"]
		if !exists {
			continue
		}
		svcVersion, err := valueObject.NewServiceVersion(svcVersionStr)
		if err != nil {
			continue
		}

		svcTypeStr, exists := svcEnvs["SVC_TYPE"]
		if !exists {
			continue
		}
		svcType, err := valueObject.NewServiceType(svcTypeStr)
		if err != nil {
			continue
		}

		svcStartupFileStr, exists := svcEnvs["SVC_STARTUP_FILE"]
		if !exists {
			svcStartupFileStr = ""
		}
		var svcStartupFilePtr *valueObject.UnixFilePath
		if svcStartupFileStr != "" {
			svcStartupFile, err := valueObject.NewUnixFilePath(svcStartupFileStr)
			if err != nil {
				continue
			}
			svcStartupFilePtr = &svcStartupFile
		}

		svcPortBindingsStr, exists := svcEnvs["SVC_PORT_BINDINGS"]
		if !exists {
			svcPortBindingsStr = ""
		}
		svcPortBindingsParts := strings.Split(svcPortBindingsStr, ",")
		svcPortBindings := []valueObject.PortBinding{}
		for _, svcPortBindingStr := range svcPortBindingsParts {
			if svcPortBindingStr == "" {
				continue
			}

			svcPortBinding, err := valueObject.NewPortBindingFromString(svcPortBindingStr)
			if err != nil {
				continue
			}
			svcPortBindings = append(svcPortBindings, svcPortBinding)
		}

		svcStatus := stoppedStatus
		if slices.Contains(runningServicesNames, svcName.String()) {
			svcStatus = runningStatus
		}

		servicesList = append(
			servicesList,
			entity.NewService(
				svcName,
				svcNature,
				svcType,
				svcVersion,
				svcCmd,
				svcStatus,
				svcStartupFilePtr,
				svcPortBindings,
			),
		)
	}

	if len(servicesList) == 0 {
		return servicesList, errors.New("GetInstalledServicesFailed")
	}

	return servicesList, nil
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
	supervisordService supervisordService,
) (valueObject.ServiceMetrics, error) {
	supervisordServiceMetrics := valueObject.ServiceMetrics{}

	cpuPercent := float64(0.0)
	memPercent := float32(0.0)

	pidProcesses, err := repo.getPpidEntireProcessFamily(supervisordService.MainPid)
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
		supervisordService.UptimeInSeconds,
		cpuPercent,
		memPercent,
	)

	return serviceMetrics, nil
}

func (repo ServicesQueryRepo) GetWithMetrics() ([]dto.ServiceWithMetrics, error) {
	servicesWithMetrics := []dto.ServiceWithMetrics{}

	servicesList, err := repo.Get()
	if err != nil {
		return servicesWithMetrics, err
	}

	supervisordServices, err := repo.getSupervisordServices()
	if err != nil {
		return servicesWithMetrics, err
	}

	for _, supervisordService := range supervisordServices {
		svcName, err := valueObject.NewServiceName(supervisordService.Name)
		if err != nil {
			continue
		}

		serviceEntityIndex := -1
		for svcIndex, serviceFromList := range servicesList {
			if serviceFromList.Name.String() != svcName.String() {
				continue
			}
			serviceEntityIndex = svcIndex
		}

		if serviceEntityIndex == -1 {
			continue
		}

		var serviceMetricsPtr *valueObject.ServiceMetrics
		if supervisordService.Status == "running" {
			serviceMetrics, err := repo.getSupervisordServiceMetrics(supervisordService)
			if err != nil {
				continue
			}
			serviceMetricsPtr = &serviceMetrics
		}

		serviceWithMetrics := dto.NewServiceWithMetrics(
			servicesList[serviceEntityIndex],
			serviceMetricsPtr,
		)
		servicesWithMetrics = append(servicesWithMetrics, serviceWithMetrics)
	}

	return servicesWithMetrics, nil
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
			valueObject.NewServiceVersionPanic("8.3"),
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

	postgresql := entity.NewInstallableService(
		valueObject.NewServiceNamePanic("postgresql"),
		valueObject.NewServiceTypePanic("database"),
		[]valueObject.ServiceVersion{
			valueObject.NewServiceVersionPanic("16"),
			valueObject.NewServiceVersionPanic("15"),
			valueObject.NewServiceVersionPanic("14"),
			valueObject.NewServiceVersionPanic("13"),
			valueObject.NewServiceVersionPanic("12"),
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
		postgresql,
		redisService,
	}, nil
}
