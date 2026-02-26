package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type O11yOverview struct {
	Hostname             tkValueObject.Fqdn               `json:"hostname"`
	UptimeSecs           uint64                           `json:"uptimeSecs"`
	UptimeRelative       tkValueObject.RelativeTime       `json:"uptimeRelative"`
	PublicIpAddress      tkValueObject.IpAddress          `json:"publicIp"`
	HardwareSpecs        valueObject.HardwareSpecs        `json:"specs"`
	CurrentResourceUsage valueObject.CurrentResourceUsage `json:"currentUsage"`
}

func NewO11yOverview(
	hostname tkValueObject.Fqdn,
	uptimeSecs uint64,
	uptimeRelative tkValueObject.RelativeTime,
	publicIpAddress tkValueObject.IpAddress,
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
