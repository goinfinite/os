package valueObject

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type StorageInfo struct {
	Total     tkValueObject.Byte `json:"total"`
	Available tkValueObject.Byte `json:"available"`
	Used      tkValueObject.Byte `json:"used"`
}

func NewStorageInfo(total, available, used tkValueObject.Byte) StorageInfo {
	return StorageInfo{
		Total:     total,
		Available: available,
		Used:      used,
	}
}
