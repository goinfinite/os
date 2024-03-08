package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type O11yOverview struct {
	Hostname             valueObject.Fqdn                 `json:"hostname"`
	UptimeSecs           uint64                           `json:"uptimeSecs"`
	PublicIpAddress      valueObject.IpAddress            `json:"publicIp"`
	HardwareSpecs        valueObject.HardwareSpecs        `json:"specs"`
	CurrentResourceUsage valueObject.CurrentResourceUsage `json:"currentUsage"`
}

func NewO11yOverview(
	hostname valueObject.Fqdn,
	uptime uint64,
	publicIpAddress valueObject.IpAddress,
	hardwareSpecs valueObject.HardwareSpecs,
	currentResourceUsage valueObject.CurrentResourceUsage,
) O11yOverview {
	return O11yOverview{
		Hostname:             hostname,
		UptimeSecs:           uptime,
		PublicIpAddress:      publicIpAddress,
		HardwareSpecs:        hardwareSpecs,
		CurrentResourceUsage: currentResourceUsage,
	}
}
