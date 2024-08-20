package valueObject

type CurrentResourceUsage struct {
	CpuUsagePercent     float64 `json:"cpuUsagePercent"`
	MemUsagePercent     float64 `json:"memUsagePercent"`
	StorageUsagePercent float64 `json:"storageUsage"`
}

func NewCurrentResourceUsage(
	cpuUsagePercent, memUsagePercent, storageUsagePercent float64,
) CurrentResourceUsage {
	return CurrentResourceUsage{
		CpuUsagePercent:     cpuUsagePercent,
		MemUsagePercent:     memUsagePercent,
		StorageUsagePercent: storageUsagePercent,
	}
}
