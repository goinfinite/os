package valueObject

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

type MktplaceItemDescription string

func NewMktplaceItemDescription(value string) (MktplaceItemDescription, error) {
	isTooShort := len(value) < 2
	isTooLong := len(value) > 512

	if isTooShort || isTooLong {
		return "", errors.New("InvalidMktItemDescription")
	}

	return MktplaceItemDescription(value), nil
}

func NewMktplaceItemDescriptionPanic(value string) MktplaceItemDescription {
	mid, err := NewMktplaceItemDescription(value)
	if err != nil {
		panic(err)
	}

	return mid
}

func (mid MktplaceItemDescription) String() string {
	return string(mid)
}

func (midPtr *MktplaceItemDescription) UnmarshalJSON(value []byte) error {
	valueStr := string(value)
	unquotedValue := strings.Trim(valueStr, "\"")

	mid, err := NewMktplaceItemDescription(unquotedValue)
	if err != nil {
		return err
	}

	*midPtr = mid
	return nil
}

func (midPtr *MktplaceItemDescription) UnmarshalYAML(value *yaml.Node) error {
	var valueStr string
	err := value.Decode(&valueStr)
	if err != nil {
		return err
	}

	mid, err := NewMktplaceItemDescription(valueStr)
	if err != nil {
		return err
	}

	*midPtr = mid
	return nil
}
