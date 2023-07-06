package o11yInfra

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type GetOverview struct {
}

func (repo GetOverview) getUptime() (uint64, error) {
	sysinfo := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(sysinfo); err != nil {
		return 0, err
	}

	return uint64(sysinfo.Uptime), nil
}

func (repo GetOverview) getCgroupLimit(file string) (int64, error) {
	fileContent, err := os.ReadFile(file)
	if err != nil {
		return 0, err
	}

	fileContentStr := strings.TrimSpace(string(fileContent))
	val, err := strconv.ParseInt(fileContentStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (repo GetOverview) getCpuQuota() (int64, error) {
	cpuQuotaFile := "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
	cpuQuota, err := repo.getCgroupLimit(cpuQuotaFile)
	if err != nil {
		return 0, errors.New("CpuQuotaFileError")
	}

	if cpuQuota == -1 {
		cpuQuota = int64(100000 * runtime.NumCPU())
	}

	return cpuQuota, nil
}

func (repo GetOverview) getMemoryLimit() (int64, error) {
	memLimitFile := "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	memLimit, err := repo.getCgroupLimit(memLimitFile)
	if err != nil {
		return 0, errors.New("MemLimitFileError")
	}

	if memLimit == 9223372036854771712 {
		var sysInfo syscall.Sysinfo_t
		err = syscall.Sysinfo(&sysInfo)
		if err != nil {
			return 0, errors.New("GetSysInfoError")
		}

		memLimit = int64(sysInfo.Totalram * uint64(sysInfo.Unit))
	}

	return memLimit, nil
}

func (repo GetOverview) getStorageInfo() (valueObject.StorageInfo, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return valueObject.StorageInfo{}, errors.New("StorageInfoError")
	}

	storageTotal := stat.Blocks * uint64(stat.Bsize)
	storageAvailable := stat.Bavail * uint64(stat.Bsize)
	storageUsed := storageTotal - storageAvailable

	return valueObject.NewStorageInfo(
		valueObject.Byte(storageTotal),
		valueObject.Byte(storageAvailable),
		valueObject.Byte(storageUsed),
	), nil
}

func (repo GetOverview) getHardwareSpecs() (valueObject.HardwareSpecs, error) {
	cmd := exec.Command(
		"awk",
		"-F:",
		"/vendor_id/{vendor=$2} /cpu MHz/{freq=$2} END{print vendor freq}",
		"/proc/cpuinfo",
	)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("GetCpuSpecsFailed: %v", err)
		return valueObject.HardwareSpecs{}, errors.New("GetCpuSpecsFailed")
	}
	trimmedOutput := strings.TrimSpace(string(output))
	if trimmedOutput == "" {
		return valueObject.HardwareSpecs{}, errors.New("EmptyCpuSpecs")
	}

	cpuInfo := strings.Split(trimmedOutput, " ")
	if len(cpuInfo) < 2 {
		return valueObject.HardwareSpecs{}, errors.New("ParseCpuSpecsFailed")
	}

	cpuModel := strings.TrimSpace(cpuInfo[0])
	cpuFrequency := strings.TrimSpace(cpuInfo[1])
	cpuFrequencyFloat, err := strconv.ParseFloat(cpuFrequency, 64)
	if err != nil {
		log.Printf("GetCpuFrequencyFailed: %v", err)
		return valueObject.HardwareSpecs{}, errors.New("GetCpuFrequencyFailed")
	}

	cpuQuota, err := repo.getCpuQuota()
	if err != nil {
		return valueObject.HardwareSpecs{}, errors.New("GetCpuQuotaFailed")
	}
	cpuCores := uint64(cpuQuota / 100000)

	memoryLimit, err := repo.getMemoryLimit()
	if err != nil {
		return valueObject.HardwareSpecs{}, errors.New("GetMemoryLimitFailed")
	}

	storageInfo, err := repo.getStorageInfo()
	if err != nil {
		return valueObject.HardwareSpecs{}, errors.New("GetStorageInfoFailed")
	}

	return valueObject.NewHardwareSpecs(
		cpuModel,
		cpuCores,
		cpuFrequencyFloat,
		valueObject.Byte(memoryLimit),
		storageInfo.Total,
	), nil
}

func (repo GetOverview) getCurrentResourceUsage() (
	valueObject.CurrentResourceUsage,
	error,
) {
	cpuUsageFile := "/sys/fs/cgroup/cpuacct/cpuacct.usage"
	startCpuUsage, err := repo.getCgroupLimit(cpuUsageFile)
	if err != nil {
		return valueObject.CurrentResourceUsage{}, errors.New("CpuStartUsageFileError")
	}
	time.Sleep(time.Second)
	endCpuUsage, err := repo.getCgroupLimit(cpuUsageFile)
	if err != nil {
		return valueObject.CurrentResourceUsage{}, errors.New("CpuEndUsageFileError")
	}

	cpuQuota, err := repo.getCpuQuota()
	if err != nil {
		return valueObject.CurrentResourceUsage{}, errors.New("GetCpuQuotaFailed")
	}
	cpuUsagePercent := float64(endCpuUsage-startCpuUsage) / 10000000 / float64(cpuQuota)

	memUsageFile := "/sys/fs/cgroup/memory/memory.usage_in_bytes"
	memUsage, err := repo.getCgroupLimit(memUsageFile)
	if err != nil {
		return valueObject.CurrentResourceUsage{}, errors.New("MemUsageFileError")
	}
	memLimit, err := repo.getMemoryLimit()
	if err != nil {
		return valueObject.CurrentResourceUsage{}, errors.New("GetMemoryLimitFailed")
	}
	memUsagePercent := float64(memUsage) / float64(memLimit) * 100

	storageInfo, err := repo.getStorageInfo()
	if err != nil {
		return valueObject.CurrentResourceUsage{}, errors.New("GetStorageInfoFailed")
	}
	storageUsagePercent := float64(storageInfo.Used.Get()) / float64(storageInfo.Total.Get()) * 100

	return valueObject.NewCurrentResourceUsage(
		cpuUsagePercent,
		memUsagePercent,
		storageUsagePercent,
	), nil
}

func (repo GetOverview) Get() (entity.O11yOverview, error) {
	hostnameStr, err := os.Hostname()
	if err != nil {
		hostnameStr = "localhost"
	}

	isVirtualHostEnvSet := os.Getenv("VIRTUAL_HOST") != ""
	if isVirtualHostEnvSet {
		hostnameStr = os.Getenv("VIRTUAL_HOST")
	}

	hostname, err := valueObject.NewFqdn(hostnameStr)
	if err != nil {
		return entity.O11yOverview{}, errors.New("GetHostnameFailed")
	}

	runtimeContext, err := infraHelper.GetRuntimeContext()
	if err != nil {
		runtimeContext, _ = valueObject.NewRuntimeContext("vm")
	}

	uptime, err := repo.getUptime()
	if err != nil {
		uptime = 0
	}

	publicIpAddress, err := infraHelper.GetPublicIpAddress()
	if err != nil {
		publicIpAddress, _ = valueObject.NewIpAddress("0.0.0.0")
	}

	hardwareSpecs, err := repo.getHardwareSpecs()
	if err != nil {
		log.Printf("GetHardwareSpecsFailed: %v", err)
		return entity.O11yOverview{}, errors.New("GetHardwareSpecsFailed")
	}

	currentResourceUsage, err := repo.getCurrentResourceUsage()
	if err != nil {
		log.Printf("GetCurrentResourceUsageFailed: %v", err)
		return entity.O11yOverview{}, errors.New("GetCurrentResourceUsageFailed")
	}

	return entity.NewO11yOverview(
		hostname,
		runtimeContext,
		uptime,
		publicIpAddress,
		hardwareSpecs,
		currentResourceUsage,
	), nil
}
