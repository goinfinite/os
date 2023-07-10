package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type Service struct {
	Name            valueObject.ServiceName   `json:"name"`
	Status          valueObject.ServiceStatus `json:"status"`
	Pid             *uint32                   `json:"pid,omitempty"`
	Uptime          *float64                  `json:"uptime,omitempty"`
	CpuUsagePercent *float64                  `json:"cpuUsagePercent,omitempty"`
	MemUsagePercent *float32                  `json:"memUsagePercent,omitempty"`
}

func NewService(
	name valueObject.ServiceName,
	status valueObject.ServiceStatus,
	pid *uint32,
	uptime *float64,
	cpuUsagePercent *float64,
	memUsagePercent *float32,
) Service {
	return Service{
		Name:            name,
		Status:          status,
		Pid:             pid,
		Uptime:          uptime,
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
	}
}
