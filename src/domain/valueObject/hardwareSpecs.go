package valueObject

type HardwareSpecs struct {
	CpuModel     string  `json:"cpuModel"`
	CpuCores     float64 `json:"cpuCores"`
	CpuFrequency float64 `json:"cpuFrequency"`
	MemoryTotal  Byte    `json:"memoryTotal"`
	StorageTotal Byte    `json:"storageTotal"`
}

func NewHardwareSpecs(
	cpuModel string,
	cpuCores, cpuFrequency float64,
	memoryTotal, storageTotal Byte,
) HardwareSpecs {
	return HardwareSpecs{
		CpuModel:     cpuModel,
		CpuCores:     cpuCores,
		CpuFrequency: cpuFrequency,
		MemoryTotal:  memoryTotal,
		StorageTotal: storageTotal,
	}
}
