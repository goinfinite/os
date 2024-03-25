package valueObject

import "errors"

type MktplaceItemDescription string

func NewMktplaceItemDescription(value string) (MktplaceItemDescription, error) {
	isTooShort := len(value) < 2
	isTooLong := len(value) > 512

	if isTooShort || isTooLong {
		return "", errors.New("InvalidMktplaceItemDescription")
	}

	return MktplaceItemDescription(value), nil
}

func NewMktplaceItemDescriptionPanic(value string) MktplaceItemDescription {
	comment, err := NewMktplaceItemDescription(value)
	if err != nil {
		panic(err)
	}

	return comment
}

func (mktplaceItemDesc MktplaceItemDescription) String() string {
	return string(mktplaceItemDesc)
}
