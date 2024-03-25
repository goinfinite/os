package valueObject

type DataField struct {
	Key   DataFieldKey   `json:"key"`
	Value DataFieldValue `json:"value"`
}

func NewDataField(
	key DataFieldKey,
	value DataFieldValue,
) DataField {
	return DataField{
		Key:   key,
		Value: value,
	}
}
