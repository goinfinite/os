package valueObject

type StorageInfo struct {
	Total     Byte `json:"total"`
	Available Byte `json:"available"`
	Used      Byte `json:"used"`
}

func NewStorageInfo(total, available, used Byte) StorageInfo {
	return StorageInfo{
		Total:     total,
		Available: available,
		Used:      used,
	}
}
