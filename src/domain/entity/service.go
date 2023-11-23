package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Service struct {
	Name            valueObject.ServiceName   `json:"name"`
	Type            valueObject.ServiceType   `json:"type"`
	Command         valueObject.UnixCommand   `json:"command"`
	Status          valueObject.ServiceStatus `json:"status"`
	Port            *valueObject.NetworkPort  `json:"port"`
	Pids            []uint32                  `json:"pids,omitempty"`
	UptimeSecs      *int64                    `json:"uptimeSecs,omitempty"`
	CpuUsagePercent *float64                  `json:"cpuUsagePercent,omitempty"`
	MemUsagePercent *float32                  `json:"memUsagePercent,omitempty"`
}

func NewService(
	name valueObject.ServiceName,
	svcType valueObject.ServiceType,
	command valueObject.UnixCommand,
	status valueObject.ServiceStatus,
	port *valueObject.NetworkPort,
	pids []uint32,
	uptimeSecs *int64,
	cpuUsagePercent *float64,
	memUsagePercent *float32,
) Service {
	return Service{
		Name:            name,
		Type:            svcType,
		Command:         command,
		Status:          status,
		Port:            port,
		Pids:            pids,
		UptimeSecs:      uptimeSecs,
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
	}
}
