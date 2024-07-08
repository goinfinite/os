package servicesInfra

import (
	"embed"
	"errors"
	"io/fs"
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

//go:embed assets/*
var assets embed.FS

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

func (repo *ServicesQueryRepo) parseManifestCmdSteps(
	stepsType string,
	rawCmdSteps interface{},
) (cmdSteps []valueObject.UnixCommand, err error) {
	cmdStepsMap, assertOk := rawCmdSteps.([]interface{})
	if !assertOk {
		return cmdSteps, errors.New("InvalidCmdSteps")
	}

	for stepIndex, rawCmd := range cmdStepsMap {
		command, err := valueObject.NewUnixCommand(rawCmd)
		if err != nil {
			log.Printf("(%sCmdSteps) [Index %d] %s", stepsType, stepIndex, err)
			return cmdSteps, err
		}
		cmdSteps = append(cmdSteps, command)
	}

	return cmdSteps, nil
}

func (repo *ServicesQueryRepo) installableServiceFactory(
	serviceFilePath valueObject.UnixFilePath,
) (installableService entity.InstallableService, err error) {
	serviceMap, err := infraHelper.EmbedSerializedDataToMap(&assets, serviceFilePath)
	if err != nil {
		return installableService, err
	}

	requiredParams := []string{
		"name", "nature", "type", "command", "description", "installCmdSteps",
	}
	for _, requiredParam := range requiredParams {
		if serviceMap[requiredParam] != nil {
			continue
		}

		return installableService, errors.New("MissingParam: " + requiredParam)
	}

	serviceName, err := valueObject.NewServiceName(serviceMap["name"])
	if err != nil {
		return installableService, err
	}
	serviceNameStr := serviceName.String()

	serviceNature, err := valueObject.NewServiceNature(serviceMap["nature"])
	if err != nil {
		return installableService, err
	}

	serviceType, err := valueObject.NewServiceType(serviceMap["type"])
	if err != nil {
		return installableService, err
	}

	serviceCommand, err := valueObject.NewUnixCommand(serviceMap["command"])
	if err != nil {
		return installableService, err
	}

	serviceDescription, err := valueObject.NewServiceDescription(serviceMap["description"])
	if err != nil {
		return installableService, err
	}

	serviceVersions := []valueObject.ServiceVersion{}
	if serviceMap["versions"] != nil {
		versionsMap, assertOk := serviceMap["versions"].([]interface{})
		if !assertOk {
			return installableService, errors.New("InvalidServiceVersions")
		}
		for versionIndex, rawVersion := range versionsMap {
			version, err := valueObject.NewServiceVersion(rawVersion)
			if err != nil {
				log.Printf("(%s) [Index %d] %s", serviceNameStr, versionIndex, err)
				continue
			}
			serviceVersions = append(serviceVersions, version)
		}
	}

	portBindings := []valueObject.PortBinding{}
	if serviceMap["portBindings"] != nil {
		portBindingsMap, assertOk := serviceMap["portBindings"].([]interface{})
		if !assertOk {
			return installableService, errors.New("InvalidPortBindings")
		}
		for portIndex, rawPortBinding := range portBindingsMap {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				log.Printf("(%s) [Index: %d] %s", serviceNameStr, portIndex, err)
				continue
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	installCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["installCmdSteps"] != nil {
		installCmdSteps, err = repo.parseManifestCmdSteps(
			"Install", serviceMap["installCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	uninstallCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["uninstallCmdSteps"] != nil {
		uninstallCmdSteps, err = repo.parseManifestCmdSteps(
			"Uninstall", serviceMap["uninstallCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	uninstallFileNames := []valueObject.UnixFileName{}
	if serviceMap["uninstallFileNames"] != nil {
		filesMap, assertOk := serviceMap["uninstallFileNames"].([]interface{})
		if !assertOk {
			return installableService, errors.New("InvalidUninstallFileNames")
		}
		for fileIndex, rawFileName := range filesMap {
			fileName, err := valueObject.NewUnixFileName(rawFileName)
			if err != nil {
				log.Printf("(%s) [Index %d] %s", serviceNameStr, fileIndex, err)
				continue
			}
			uninstallFileNames = append(uninstallFileNames, fileName)
		}
	}

	preStartCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["preStartCmdSteps"] != nil {
		preStartCmdSteps, err = repo.parseManifestCmdSteps(
			"PreStart", serviceMap["preStartCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	postStartCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["postStartCmdSteps"] != nil {
		postStartCmdSteps, err = repo.parseManifestCmdSteps(
			"PostStart", serviceMap["postStartCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	preStopCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["preStopCmdSteps"] != nil {
		preStopCmdSteps, err = repo.parseManifestCmdSteps(
			"PreStop", serviceMap["preStopCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	postStopCmdSteps := []valueObject.UnixCommand{}
	if serviceMap["postStopCmdSteps"] != nil {
		postStopCmdSteps, err = repo.parseManifestCmdSteps(
			"PostStop", serviceMap["postStopCmdSteps"],
		)
		if err != nil {
			return installableService, err
		}
	}

	var startupFilePtr *valueObject.UnixFilePath
	if serviceMap["startupFile"] != nil {
		startupFile, err := valueObject.NewUnixFilePath(serviceMap["startupFile"])
		if err != nil {
			return installableService, err
		}
		startupFilePtr = &startupFile
	}

	var estimatedSizeBytesPtr *valueObject.Byte
	if serviceMap["estimatedSizeBytes"] == nil {
		estimatedSizeBytes, err := valueObject.NewByte(serviceMap["estimatedSizeBytes"])
		if err != nil {
			return installableService, err
		}
		estimatedSizeBytesPtr = &estimatedSizeBytes
	}

	var avatarUrlPtr *valueObject.Url
	if serviceMap["avatarUrl"] != nil {
		avatarUrl, err := valueObject.NewUrl(serviceMap["avatarUrl"])
		if err != nil {
			return installableService, err
		}
		avatarUrlPtr = &avatarUrl
	}

	return entity.NewInstallableService(
		serviceName, serviceNature, serviceType, serviceCommand,
		serviceDescription, serviceVersions, portBindings, installCmdSteps,
		uninstallCmdSteps, uninstallFileNames, preStartCmdSteps, postStartCmdSteps,
		preStopCmdSteps, postStopCmdSteps, startupFilePtr, estimatedSizeBytesPtr,
		avatarUrlPtr,
	), nil
}

func (repo *ServicesQueryRepo) ReadInstallables() ([]entity.InstallableService, error) {
	installableServices := []entity.InstallableService{}

	serviceFiles, err := fs.ReadDir(assets, "assets")
	if err != nil {
		return installableServices, errors.New("ReadServiceFilesError: " + err.Error())
	}

	for _, serviceFile := range serviceFiles {
		serviceFileName := serviceFile.Name()
		rawServiceFilePath := "assets/" + serviceFileName
		serviceFilePath, err := valueObject.NewUnixFilePath(rawServiceFilePath)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawServiceFilePath)
			continue
		}
		serviceFilePathStr := serviceFilePath.String()

		installableService, err := repo.installableServiceFactory(serviceFilePath)
		if err != nil {
			log.Printf(
				"ReadServiceFileError (%s): %s", serviceFilePathStr, err.Error(),
			)
			continue
		}

		installableServices = append(installableServices, installableService)
	}

	return installableServices, nil
}
