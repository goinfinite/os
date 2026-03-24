package valueObject

import (
	"fmt"
	"strings"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type HardwareSpecs struct {
	CpuModel     string  `json:"cpuModel"`
	CpuCores     float64 `json:"cpuCores"`
	CpuFrequency float64 `json:"cpuFrequency"`
	MemoryTotal  tkValueObject.Byte `json:"memoryTotal"`
	StorageTotal tkValueObject.Byte `json:"storageTotal"`
}

func NewHardwareSpecs(
	cpuModel string,
	cpuCores, cpuFrequency float64,
	memoryTotal, storageTotal tkValueObject.Byte,
) HardwareSpecs {
	return HardwareSpecs{
		CpuModel:     cpuModel,
		CpuCores:     cpuCores,
		CpuFrequency: cpuFrequency,
		MemoryTotal:  memoryTotal,
		StorageTotal: storageTotal,
	}
}

func (vo HardwareSpecs) String() string {
	cpuModelNameParts := strings.Split(vo.CpuModel, " ")
	if len(cpuModelNameParts) > 4 {
		cpuModelNameParts = cpuModelNameParts[:4]
	}
	cpuModelNameStr := strings.Join(cpuModelNameParts, " ")

	cpuFrequencyGhz := vo.CpuFrequency / 1000

	return fmt.Sprintf(
		"%s (%.0fc@%.1f GHz) ‖ %s RAM",
		cpuModelNameStr, vo.CpuCores,
		cpuFrequencyGhz, vo.MemoryTotal.StringWithSuffix(),
	)
}
