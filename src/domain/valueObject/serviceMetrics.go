package valueObject

type ServiceMetrics struct {
	Pids            []uint32 `json:"pids"`
	UptimeSecs      int64    `json:"uptimeSecs"`
	CpuUsagePercent float64  `json:"cpuUsagePercent"`
	MemUsagePercent float32  `json:"memUsagePercent"`
}

func NewServiceMetrics(
	pids []uint32,
	uptimeSecs int64,
	cpuUsagePercent float64,
	memUsagePercent float32,
) ServiceMetrics {
	return ServiceMetrics{
		Pids:            pids,
		UptimeSecs:      uptimeSecs,
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
	}
}
