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
	cpuCores float64,
	cpuFrequency float64,
	memoryTotal Byte,
	storageTotal Byte,
) HardwareSpecs {
	return HardwareSpecs{
		CpuModel:     cpuModel,
		CpuCores:     cpuCores,
		CpuFrequency: cpuFrequency,
		MemoryTotal:  memoryTotal,
		StorageTotal: storageTotal,
	}
}
