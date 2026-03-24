package valueObject

import (
	"errors"
	"strings"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type PortBinding struct {
	Port     tkValueObject.NetworkPort     `json:"port"`
	Protocol tkValueObject.NetworkProtocol `json:"protocol"`
}

func NewPortBinding(value interface{}) (portBinding PortBinding, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
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

	port, err := tkValueObject.NewNetworkPort(bindingParts[0])
	if err != nil {
		return portBinding, err
	}

	protocol, err := tkValueObject.NewNetworkProtocol(bindingParts[1])
	if err != nil {
		return portBinding, err
	}

	return PortBinding{
		Port: port, Protocol: protocol,
	}, nil
}

func (vo PortBinding) GetPort() tkValueObject.NetworkPort {
	return vo.Port
}

func (vo PortBinding) GetProtocol() tkValueObject.NetworkProtocol {
	return vo.Protocol
}

func (vo PortBinding) String() string {
	return vo.Port.String() + "/" + vo.Protocol.String()
}
