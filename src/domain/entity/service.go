package entity

import "github.com/speedianet/os/src/domain/valueObject"

type Service struct {
	Name            valueObject.ServiceName    `json:"name"`
	Nature          valueObject.ServiceNature  `json:"nature"`
	Type            valueObject.ServiceType    `json:"type"`
	Version         valueObject.ServiceVersion `json:"version"`
	Command         valueObject.UnixCommand    `json:"command"`
	Status          valueObject.ServiceStatus  `json:"status"`
	StartupFile     *valueObject.UnixFilePath  `json:"startupFile,omitempty"`
	Ports           []valueObject.NetworkPort  `json:"ports,omitempty"`
	Pids            []uint32                   `json:"pids,omitempty"`
	UptimeSecs      *int64                     `json:"uptimeSecs,omitempty"`
	CpuUsagePercent *float64                   `json:"cpuUsagePercent,omitempty"`
	MemUsagePercent *float32                   `json:"memUsagePercent,omitempty"`
}

func NewService(
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	svcType valueObject.ServiceType,
	version valueObject.ServiceVersion,
	command valueObject.UnixCommand,
	status valueObject.ServiceStatus,
	startupFile *valueObject.UnixFilePath,
	ports []valueObject.NetworkPort,
	pids []uint32,
	uptimeSecs *int64,
	cpuUsagePercent *float64,
	memUsagePercent *float32,
) Service {
	return Service{
		Name:            name,
		Nature:          nature,
		Type:            svcType,
		Version:         version,
		Command:         command,
		Status:          status,
		StartupFile:     startupFile,
		Ports:           ports,
		Pids:            pids,
		UptimeSecs:      uptimeSecs,
		CpuUsagePercent: cpuUsagePercent,
		MemUsagePercent: memUsagePercent,
	}
}
