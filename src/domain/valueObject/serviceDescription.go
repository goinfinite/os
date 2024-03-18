package valueObject

import "errors"

type ServiceDescription string

func NewServiceDescription(value string) (ServiceDescription, error) {
	isTooShort := len(value) < 2
	isTooLong := len(value) > 512

	if isTooShort || isTooLong {
		return "", errors.New("InvalidServiceDescription")
	}

	return ServiceDescription(value), nil
}

func NewServiceDescriptionPanic(value string) ServiceDescription {
	comment, err := NewServiceDescription(value)
	if err != nil {
		panic(err)
	}

	return comment
}

func (svcDesc ServiceDescription) String() string {
	return string(svcDesc)
}
