package valueObject

import (
	"errors"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type PortBinding struct {
	Port     NetworkPort     `json:"port"`
	Protocol NetworkProtocol `json:"protocol"`
}

func NewPortBinding(value interface{}) (portBinding PortBinding, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return portBinding, errors.New("PortBindingValueMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if len(stringValue) == 0 {
		return portBinding, errors.New("EmptyPortBinding")
	}

	if !strings.Contains(stringValue, "/") {
		stringValue += "/tcp"
	}

	bindingParts := strings.Split(stringValue, "/")
	if len(bindingParts) != 2 {
		return portBinding, errors.New("InvalidPortBinding")
	}

	port, err := NewNetworkPort(bindingParts[0])
	if err != nil {
		return portBinding, err
	}

	protocol, err := NewNetworkProtocol(bindingParts[1])
	if err != nil {
		return portBinding, err
	}

	return PortBinding{
		Port: port, Protocol: protocol,
	}, nil
}

func (vo PortBinding) GetPort() NetworkPort {
	return vo.Port
}

func (vo PortBinding) GetProtocol() NetworkProtocol {
	return vo.Protocol
}

func (vo PortBinding) String() string {
	return vo.Port.String() + "/" + vo.Protocol.String()
}
