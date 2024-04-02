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

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

const PublicIpTransientKey string = "PublicIp"

type O11yQueryRepo struct {
	transientDbSvc *internalDbInfra.TransientDatabaseService
}

func NewO11yQueryRepo(
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *O11yQueryRepo {
	return &O11yQueryRepo{
		transientDbSvc: transientDbSvc,
	}
}

func (repo *O11yQueryRepo) getUptime() (uint64, error) {
	sysinfo := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(sysinfo); err != nil {
		return 0, err
	}

	return uint64(sysinfo.Uptime), nil
}

func (repo *O11yQueryRepo) getPublicIpAddress() (valueObject.IpAddress, error) {
	var ipAddress valueObject.IpAddress

	cachedIpAddressStr, err := repo.transientDbSvc.Get(PublicIpTransientKey)
	if err == nil {
		return valueObject.NewIpAddress(cachedIpAddressStr)
	}

	rawIpEntry, err := infraHelper.RunCmd(
		"dig", "+short", "TXT", "o-o.myaddr.l.google.com", "@ns1.google.com",
	)
	if err != nil {
		rawIpEntry, err = infraHelper.RunCmd(
			"dig", "+short", "TXT", "CH", "whoami.cloudflare", "@1.1.1.1",
		)
		if err != nil {
			return ipAddress, errors.New("GetPublicIpFailed: " + err.Error())
		}
	}

	rawIpEntry = strings.Trim(rawIpEntry, `"`)
	if rawIpEntry == "" {
		return ipAddress, errors.New("GetPublicIpFailed: NoIpEntry")
	}

	ipAddress, err = valueObject.NewIpAddress(rawIpEntry)
	if err != nil {
		return ipAddress, err
	}

	err = repo.transientDbSvc.Set(PublicIpTransientKey, ipAddress.String())
	if err != nil {
		return ipAddress, errors.New("PersistPublicIpFailed: " + err.Error())
	}

	return ipAddress, nil
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

func (repo *O11yQueryRepo) getCpuQuota() (int64, error) {
	cpuQuotaFile := "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
	if repo.isCgroupV2() {
		cpuQuotaFile = "/sys/fs/cgroup/cpu.max"
	}

	cpuQuotaStr, err := repo.getFileContent(cpuQuotaFile)
	if err != nil {
		cpuQuotaStr = "max"
	}
	if repo.isCgroupV2() {
		cpuQuotaStr = strings.Split(cpuQuotaStr, " ")[0]
	}

	cpuQuotaInt, err := strconv.ParseInt(cpuQuotaStr, 10, 64)
	if err != nil || cpuQuotaStr == "max" || cpuQuotaStr == "-1" {
		cpuQuotaInt = int64(100000 * runtime.NumCPU())
	}

	return cpuQuotaInt, nil
}

func (repo *O11yQueryRepo) getMemoryLimit() (int64, error) {
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

	return memLimitInt, nil
}

func (repo *O11yQueryRepo) getStorageInfo() (valueObject.StorageInfo, error) {
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

	cpuQuotaUs, err := repo.getCpuQuota()
	if err != nil {
		return 0, errors.New("GetCpuQuotaFailed")
	}

	cpuUsageUs := float64(endCpuUsageInt - startCpuUsageInt)
	if !repo.isCgroupV2() {
		cpuUsageUs = cpuUsageUs / 1000
	}
	cpuUsagePercent := cpuUsageUs / float64(cpuQuotaUs) * 100

	return cpuUsagePercent, nil
}

func (repo *O11yQueryRepo) getMemUsagePercent() (float64, error) {
	memUsageFile := "/sys/fs/cgroup/memory/memory.usage_in_bytes"
	if repo.isCgroupV2() {
		memUsageFile = "/sys/fs/cgroup/memory.current"
	}

	memUsage, err := repo.getFileContent(memUsageFile)
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

		memUsage = strings.TrimSpace(string(cmdOutput))
	}
	memUsageFloat, err := strconv.ParseFloat(memUsage, 64)
	if err != nil {
		return 0, errors.New("ParseMemUsageFailed")
	}

	memLimit, err := repo.getMemoryLimit()
	if err != nil {
		return 0, errors.New("GetMemoryLimitFailed")
	}
	memUsagePercent := memUsageFloat / float64(memLimit) * 100

	return memUsagePercent, nil
}

func (repo *O11yQueryRepo) getCurrentResourceUsage() (
	valueObject.CurrentResourceUsage,
	error,
) {
	cpuUsagePercent, err := repo.getCpuUsagePercent()
	if err != nil {
		return valueObject.CurrentResourceUsage{}, err
	}
	memUsagePercent, err := repo.getMemUsagePercent()
	if err != nil {
		return valueObject.CurrentResourceUsage{}, err
	}

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

func (repo *O11yQueryRepo) GetOverview() (entity.O11yOverview, error) {
	var o11yOverview entity.O11yOverview

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

	uptime, err := repo.getUptime()
	if err != nil {
		uptime = 0
	}

	publicIpAddress, err := repo.getPublicIpAddress()
	if err != nil {
		log.Printf("GetPublicIpAddressError: %s", err.Error())
		publicIpAddress, _ = valueObject.NewIpAddress("0.0.0.0")
	}

	hardwareSpecs, err := repo.getHardwareSpecs()
	if err != nil {
		return o11yOverview, errors.New("GetHardwareSpecsFailed: " + err.Error())
	}

	currentResourceUsage, err := repo.getCurrentResourceUsage()
	if err != nil {
		return o11yOverview, errors.New("GetCurrentResourceUsageFailed: " + err.Error())
	}

	return entity.NewO11yOverview(
		hostname,
		uptime,
		publicIpAddress,
		hardwareSpecs,
		currentResourceUsage,
	), nil
}
