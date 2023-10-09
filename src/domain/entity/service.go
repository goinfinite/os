package entity

import "github.com/speedianet/sam/src/domain/valueObject"

type Service struct {
	Name            valueObject.ServiceName   `json:"name"`
	Status          valueObject.ServiceStatus `json:"status"`
	Pids            *[]uint32                 `json:"pids,omitempty"`
	UptimeSecs      *int64                    `json:"uptimeSecs,omitempty"`
	CpuUsagePercent *float64                  `json:"cpuUsagePercent,omitempty"`
	MemUsagePercent *float32                  `json:"memUsagePercent,omitempty"`
}

func NewService(
	name valueObject.ServiceName,
	status valueObject.ServiceStatus,
	pids *[]uint32,
	uptimeSecs *int64,
	cpuUsagePercent *float64,
	memUsagePercent *float32,
) Service {
	return Service{
		Name:            name,
		Status:          status,
		Pids:            pids,
		UptimeSecs:      uptimeSecs,
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
	}
}
