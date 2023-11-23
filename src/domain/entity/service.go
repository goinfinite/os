package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Service struct {
	Name            valueObject.ServiceName   `json:"name"`
	Type            valueObject.ServiceType   `json:"type"`
	Status          valueObject.ServiceStatus `json:"status"`
	Command         *valueObject.UnixCommand  `json:"command,omitempty"`
	Port            *valueObject.NetworkPort  `json:"port,omitempty"`
	Pids            []uint32                  `json:"pids,omitempty"`
	UptimeSecs      *int64                    `json:"uptimeSecs,omitempty"`
	CpuUsagePercent *float64                  `json:"cpuUsagePercent,omitempty"`
	MemUsagePercent *float32                  `json:"memUsagePercent,omitempty"`
}

func NewService(
	name valueObject.ServiceName,
	svcType valueObject.ServiceType,
	status valueObject.ServiceStatus,
	command *valueObject.UnixCommand,
	port *valueObject.NetworkPort,
	pids []uint32,
	uptimeSecs *int64,
	cpuUsagePercent *float64,
	memUsagePercent *float32,
) Service {
	return Service{
		Name:            name,
		Type:            svcType,
		Status:          status,
		Command:         command,
		Port:            port,
		Pids:            pids,
		UptimeSecs:      uptimeSecs,
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
	}
}
