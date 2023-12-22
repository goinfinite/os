package valueObject

import (
	"errors"
	"strings"
)

type PortBinding struct {
	Port     NetworkPort     `json:"port"`
	Protocol NetworkProtocol `json:"protocol"`
}

func NewPortBinding(
	port NetworkPort,
	protocol NetworkProtocol,
) PortBinding {
	return PortBinding{
		Port:     port,
		Protocol: protocol,
	}
}

func NewPortBindingFromString(value string) (PortBinding, error) {
	var portBinding PortBinding

	if value == "" {
		return portBinding, errors.New("InvalidPortBinding")
	}

	if !strings.Contains(value, "/") {
		return portBinding, errors.New("InvalidPortBinding")
	}

	specParts := strings.Split(value, "/")
	if len(specParts) != 2 {
		return portBinding, errors.New("InvalidPortBinding")
	}

	port, err := NewNetworkPort(specParts[0])
	if err != nil {
		return portBinding, err
	}

	protocol, err := NewNetworkProtocol(specParts[1])
	if err != nil {
		return portBinding, err
	}

	return NewPortBinding(port, protocol), nil
}

func (portBinding PortBinding) GetPort() NetworkPort {
	return portBinding.Port
}

func (portBinding PortBinding) GetProtocol() NetworkProtocol {
	return portBinding.Protocol
}

func (portBinding PortBinding) String() string {
	return portBinding.Port.String() + "/" + portBinding.Protocol.String()
}
