package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type O11yOverview struct {
	Hostname             valueObject.Fqdn                 `json:"hostname"`
	UptimeSecs           uint64                           `json:"uptimeSecs"`
	UptimeRelative       valueObject.RelativeTime         `json:"uptimeRelative"`
	PublicIpAddress      valueObject.IpAddress            `json:"publicIp"`
	HardwareSpecs        valueObject.HardwareSpecs        `json:"specs"`
	CurrentResourceUsage valueObject.CurrentResourceUsage `json:"currentUsage"`
}

func NewO11yOverview(
	hostname valueObject.Fqdn,
	uptimeSecs uint64,
	uptimeRelative valueObject.RelativeTime,
	publicIpAddress valueObject.IpAddress,
	hardwareSpecs valueObject.HardwareSpecs,
	currentResourceUsage valueObject.CurrentResourceUsage,
) O11yOverview {
	return O11yOverview{
		Hostname:             hostname,
		UptimeSecs:           uptimeSecs,
		UptimeRelative:       uptimeRelative,
		PublicIpAddress:      publicIpAddress,
		HardwareSpecs:        hardwareSpecs,
		CurrentResourceUsage: currentResourceUsage,
	}
}
