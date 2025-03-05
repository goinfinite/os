package o11yInfra

import (
	"errors"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

const PublicIpTransientKey string = "PublicIp"

type O11yQueryRepo struct {
	transientDbSvc *internalDbInfra.TransientDatabaseService
}

func NewO11yQueryRepo(
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *O11yQueryRepo {
	return &O11yQueryRepo{transientDbSvc: transientDbSvc}
}

func (repo *O11yQueryRepo) getUptime() (uint64, error) {
	nowEpoch := valueObject.NewUnixTimeNow()
	rawFirstPidEpoch, err := infraHelper.RunCmd(infraHelper.RunCmdConfigs{
		Command:               "stat -c '%Y' /proc/1",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return 0, errors.New("ReadFirstPidEpochFailed")
	}
	firstPidEpoch, err := valueObject.NewUnixTime(rawFirstPidEpoch)
	if err != nil {
		return 0, errors.New("ParseFirstPidEpochFailed")
	}

	uptimeSecs := nowEpoch.Int64() - firstPidEpoch.Int64()
	return uint64(uptimeSecs), nil
}

func (repo *O11yQueryRepo) ReadServerPublicIpAddress() (
	ipAddress valueObject.IpAddress,
	err error,
) {
	cachedIpAddressStr, err := repo.transientDbSvc.Get(PublicIpTransientKey)
	if err == nil {
		return valueObject.NewIpAddress(cachedIpAddressStr)
	}

	serverPublicIpAddress, err := infraHelper.ReadServerPublicIpAddress()
	if err != nil {
		return ipAddress, errors.New("ReadServerPublicIpAddressError: " + err.Error())
	}

	err = repo.transientDbSvc.Set(PublicIpTransientKey, serverPublicIpAddress.String())
	if err != nil {
		return ipAddress, errors.New("PersistPublicIpFailed: " + err.Error())
	}

	return serverPublicIpAddress, nil
}

func (repo *O11yQueryRepo) isCgroupV2() bool {
	_, err := os.Stat("/sys/fs/cgroup/cpu.max")
	return err == nil
}

func (repo *O11yQueryRepo) getFileContent(file string) (string, error) {
	fileContent, err := infraHelper.GetFileContent(file)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(fileContent), nil
}

func (repo *O11yQueryRepo) getCpuCores() (float64, error) {
	cpuQuotaStr, err := repo.getFileContent("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	if err != nil {
		cpuQuotaStr = "max"
	}

	cpuPeriodStr, err := repo.getFileContent("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
	if err != nil {
		cpuPeriodStr = "100000"
	}

	if repo.isCgroupV2() {
		cpuQuotaPeriodStr, err := repo.getFileContent("/sys/fs/cgroup/cpu.max")
		if err != nil {
			cpuQuotaPeriodStr = "max 100000"
		}
		cpuQuotaPeriodParts := strings.Split(cpuQuotaPeriodStr, " ")
		if len(cpuQuotaPeriodParts) > 1 {
			cpuQuotaStr = cpuQuotaPeriodParts[0]
			cpuPeriodStr = cpuQuotaPeriodParts[1]
		}
	}

	cpuQuotaInt, err := strconv.ParseFloat(cpuQuotaStr, 64)
	if err != nil || cpuQuotaStr == "max" || cpuQuotaStr == "-1" {
		cpuQuotaInt = float64(100000 * runtime.NumCPU())
	}

	cpuPeriodInt, err := strconv.ParseFloat(cpuPeriodStr, 64)
	if err != nil {
		cpuPeriodInt = 100000
	}

	return cpuQuotaInt / cpuPeriodInt, nil
}

func (repo *O11yQueryRepo) getMemoryLimit() (valueObject.Byte, error) {
	memLimitFile := "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	if repo.isCgroupV2() {
		memLimitFile = "/sys/fs/cgroup/memory.max"
	}

	memLimit, err := repo.getFileContent(memLimitFile)
	if err != nil {
		memLimit = "max"
	}

	memLimitInt, err := strconv.ParseInt(memLimit, 10, 64)
	if err != nil || memLimit == "9223372036854771712" || memLimit == "max" {
		var sysInfo syscall.Sysinfo_t
		err = syscall.Sysinfo(&sysInfo)
		if err != nil {
			return 0, errors.New("GetSysInfoError")
		}

		memLimitInt = int64(sysInfo.Totalram * uint64(sysInfo.Unit))
	}

	return valueObject.NewByte(memLimitInt)
}

func (repo *O11yQueryRepo) getStorageInfo() (valueObject.StorageInfo, error) {
	var storageInfo valueObject.StorageInfo

	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return storageInfo, errors.New("StorageInfoError")
	}

	storageTotalUint := stat.Blocks * uint64(stat.Bsize)
	storageTotal, err := valueObject.NewByte(storageTotalUint)
	if err != nil {
		return storageInfo, err
	}

	storageAvailableUint := stat.Bavail * uint64(stat.Bsize)
	storageAvailable, err := valueObject.NewByte(storageAvailableUint)
	if err != nil {
		return storageInfo, err
	}

	storageUsedUint := storageTotalUint - storageAvailableUint
	storageUsed, err := valueObject.NewByte(storageUsedUint)
	if err != nil {
		return storageInfo, err
	}

	return valueObject.NewStorageInfo(
		storageTotal,
		storageAvailable,
		storageUsed,
	), nil
}

func (repo *O11yQueryRepo) getHardwareSpecs() (valueObject.HardwareSpecs, error) {
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

	cpuCores, err := repo.getCpuCores()
	if err != nil {
		return valueObject.HardwareSpecs{}, errors.New("GetCpuQuotaFailed")
	}

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
		memoryLimit,
		storageInfo.Total,
	), nil
}

func (repo *O11yQueryRepo) getCpuUsagePercent() (float64, error) {
	cpuUsageFile := "/sys/fs/cgroup/cpuacct/cpuacct.usage"
	if repo.isCgroupV2() {
		cpuUsageFile = "/sys/fs/cgroup/cpu.stat"
	}

	readUsageFileErr := false
	startCpuUsage, err := repo.getFileContent(cpuUsageFile)
	if err != nil {
		readUsageFileErr = true
		startCpuUsage, err = repo.getFileContent("/proc/stat")
		if err != nil {
			return 0, errors.New("CpuStartUsageFileError")
		}
		startCpuUsage = strings.Fields(startCpuUsage)[2]
	}
	time.Sleep(time.Second)
	endCpuUsage, err := repo.getFileContent(cpuUsageFile)
	if err != nil {
		readUsageFileErr = true
		endCpuUsage, err = repo.getFileContent("/proc/stat")
		if err != nil {
			return 0, errors.New("CpuEndUsageFileError")
		}
		endCpuUsage = strings.Fields(endCpuUsage)[2]
	}

	if repo.isCgroupV2() && !readUsageFileErr {
		startCpuUsage = strings.Fields(startCpuUsage)[1]
		endCpuUsage = strings.Fields(endCpuUsage)[1]
	}

	startCpuUsageInt, err := strconv.ParseInt(startCpuUsage, 10, 64)
	if err != nil {
		return 0, errors.New("ParseCpuStartUsageFailed")
	}
	endCpuUsageInt, err := strconv.ParseInt(endCpuUsage, 10, 64)
	if err != nil {
		return 0, errors.New("ParseCpuEndUsageFailed")
	}

	cpuCores, err := repo.getCpuCores()
	if err != nil {
		return 0, errors.New("GetCpuCoresFailed")
	}
	cpuCoresUs := cpuCores * 1000000

	cpuUsageUs := float64(endCpuUsageInt - startCpuUsageInt)
	if !repo.isCgroupV2() {
		cpuUsageUs = cpuUsageUs / 1000
	}
	cpuUsagePercent := (cpuUsageUs / cpuCoresUs) * 100
	if cpuUsagePercent > 100 {
		cpuUsagePercent = 100
	}

	return cpuUsagePercent, nil
}

func (repo *O11yQueryRepo) getMemUsagePercent() (float64, error) {
	memUsageFile := "/sys/fs/cgroup/memory/memory.usage_in_bytes"
	if repo.isCgroupV2() {
		memUsageFile = "/sys/fs/cgroup/memory.current"
	}

	memUsageStr, err := repo.getFileContent(memUsageFile)
	if err != nil {
		memUsageCmd := exec.Command(
			"awk",
			"/^MemTotal:/ {total=$2} /^MemAvailable:/ {available=$2} END {used=(total-available)*1024; printf \"%d\", used}",
			"/proc/meminfo",
		)
		cmdOutput, err := memUsageCmd.Output()
		if err != nil {
			return 0, errors.New("GetMemUsageFailed")
		}

		memUsageStr = strings.TrimSpace(string(cmdOutput))
	}
	memUsageFloat, err := strconv.ParseFloat(memUsageStr, 64)
	if err != nil {
		return 0, errors.New("ParseMemUsageFailed")
	}

	memLimit, err := repo.getMemoryLimit()
	if err != nil {
		return 0, errors.New("GetMemoryLimitFailed")
	}
	memUsagePercent := (memUsageFloat / float64(memLimit)) * 100
	if memUsagePercent > 100 {
		memUsagePercent = 100
	}

	return memUsagePercent, nil
}

func (repo *O11yQueryRepo) getCurrentResourceUsage() (
	resourceUsage valueObject.CurrentResourceUsage,
	err error,
) {
	cpuUsagePercent, err := repo.getCpuUsagePercent()
	if err != nil {
		return resourceUsage, err
	}
	memUsagePercent, err := repo.getMemUsagePercent()
	if err != nil {
		return resourceUsage, err
	}

	storageInfo, err := repo.getStorageInfo()
	if err != nil {
		return resourceUsage, errors.New("ReadStorageInfoFailed")
	}
	storageUsagePercent := float64(storageInfo.Used.Int64()) / float64(storageInfo.Total.Int64()) * 100

	cpuUsagePercentStr := strconv.FormatFloat(cpuUsagePercent, 'f', 0, 64)
	memUsagePercentStr := strconv.FormatFloat(memUsagePercent, 'f', 0, 64)
	storageUsagePercentStr := strconv.FormatFloat(storageUsagePercent, 'f', 0, 64)

	return valueObject.NewCurrentResourceUsage(
		cpuUsagePercent, cpuUsagePercentStr,
		memUsagePercent, memUsagePercentStr,
		storageUsagePercent, storageUsagePercentStr,
	), nil
}

func (repo *O11yQueryRepo) ReadOverview(
	withResourceUsage bool,
) (o11yOverview entity.O11yOverview, err error) {
	hostnameStr, err := os.Hostname()
	if err != nil {
		hostnameStr = "localhost"
	}

	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err == nil {
		hostnameStr = primaryVhost.String()
	}

	hostname, err := valueObject.NewFqdn(hostnameStr)
	if err != nil {
		return o11yOverview, errors.New("GetHostnameFailed")
	}

	uptimeSecs, err := repo.getUptime()
	if err != nil {
		uptimeSecs = 0
	}

	uptimeSecsDuration := time.Duration(uptimeSecs) * time.Second
	humanizedUptime := humanize.Time(time.Now().Add(-uptimeSecsDuration))
	uptimeRelative, err := valueObject.NewRelativeTime(humanizedUptime)
	if err != nil {
		uptimeRelative, _ = valueObject.NewRelativeTime("0 seconds ago")
	}

	publicIpAddress, err := repo.ReadServerPublicIpAddress()
	if err != nil {
		slog.Debug("ReadServerPublicIpAddressError", slog.String("err", err.Error()))
		publicIpAddress, _ = valueObject.NewIpAddress("0.0.0.0")
	}

	hardwareSpecs, err := repo.getHardwareSpecs()
	if err != nil {
		return o11yOverview, errors.New("ReadHardwareSpecsFailed: " + err.Error())
	}

	resourceUsage := valueObject.CurrentResourceUsage{}
	if withResourceUsage {
		currentResourceUsage, err := repo.getCurrentResourceUsage()
		if err != nil {
			return o11yOverview, errors.New("ReadCurrentResourceUsageFailed: " + err.Error())
		}
		resourceUsage = currentResourceUsage
	}

	return entity.NewO11yOverview(
		hostname, uptimeSecs, uptimeRelative, publicIpAddress, hardwareSpecs, resourceUsage,
	), nil
}
