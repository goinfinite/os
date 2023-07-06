package valueObject

type StorageInfo struct {
	Total     Byte `json:"total"`
	Available Byte `json:"available"`
	Used      Byte `json:"used"`
}

func NewStorageInfo(
	total Byte,
	available Byte,
	used Byte,
) StorageInfo {
	return StorageInfo{
		Total:     total,
		Available: available,
		Used:      used,
	}
}
