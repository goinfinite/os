package valueObject

type CurrentResourceUsage struct {
	CpuUsagePercent        float64 `json:"cpuUsagePercent"`
	CpuUsagePercentStr     string  `json:"cpuUsagePercentStr"`
	MemUsagePercent        float64 `json:"memUsagePercent"`
	MemUsagePercentStr     string  `json:"memUsagePercentStr"`
	StorageUsagePercent    float64 `json:"storageUsage"`
	StorageUsagePercentStr string  `json:"storageUsagePercentStr"`
}

func NewCurrentResourceUsage(
	cpuUsagePercent float64,
	cpuUsagePercentStr string,
	memUsagePercent float64,
	memUsagePercentStr string,
	storageUsagePercent float64,
	storageUsagePercentStr string,
) CurrentResourceUsage {
	return CurrentResourceUsage{
		CpuUsagePercent:        cpuUsagePercent,
		CpuUsagePercentStr:     cpuUsagePercentStr,
		MemUsagePercent:        memUsagePercent,
		MemUsagePercentStr:     memUsagePercentStr,
		StorageUsagePercent:    storageUsagePercent,
		StorageUsagePercentStr: storageUsagePercentStr,
	}
}
