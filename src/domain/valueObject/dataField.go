package valueObject

import (
	"gopkg.in/yaml.v3"
)

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

func (dfPtr *DataField) UnmarshalYAML(value *yaml.Node) error {
	var valuesMap map[string]string
	err := value.Decode(&valuesMap)
	if err != nil {
		return err
	}

	dfKey, err := NewDataFieldKey(valuesMap["key"])
	if err != nil {
		return err
	}

	dfValue, err := NewDataFieldValue(valuesMap["value"])
	if err != nil {
		return err
	}

	*dfPtr = NewDataField(dfKey, dfValue)

	return nil
}
