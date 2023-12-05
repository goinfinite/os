package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Service struct {
	Name            valueObject.ServiceName    `json:"name"`
	Type            valueObject.ServiceType    `json:"type"`
	Version         valueObject.ServiceVersion `json:"version"`
	Status          valueObject.ServiceStatus  `json:"status"`
	Command         *valueObject.UnixCommand   `json:"command,omitempty"`
	Ports           []valueObject.NetworkPort  `json:"ports,omitempty"`
	Pids            []uint32                   `json:"pids,omitempty"`
	UptimeSecs      *int64                     `json:"uptimeSecs,omitempty"`
	CpuUsagePercent *float64                   `json:"cpuUsagePercent,omitempty"`
	MemUsagePercent *float32                   `json:"memUsagePercent,omitempty"`
}

func NewService(
	name valueObject.ServiceName,
	svcType valueObject.ServiceType,
	version valueObject.ServiceVersion,
	status valueObject.ServiceStatus,
	command *valueObject.UnixCommand,
	ports []valueObject.NetworkPort,
	pids []uint32,
	uptimeSecs *int64,
	cpuUsagePercent *float64,
	memUsagePercent *float32,
) Service {
	return Service{
		Name:            name,
		Type:            svcType,
		Version:         version,
		Status:          status,
		Command:         command,
		Ports:           ports,
		Pids:            pids,
		UptimeSecs:      uptimeSecs,
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
	}
}
