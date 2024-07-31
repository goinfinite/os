package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type NetworkProtocol string

var ValidNetworkProtocols = []string{
	"http", "https", "ws", "wss", "grpc", "grpcs", "tcp", "udp",
}

func NewNetworkProtocol(value interface{}) (
	networkProtocol NetworkProtocol, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return networkProtocol, errors.New("NetworkProtocolMustBeString")
	}
	stringValue = strings.TrimSpace(stringValue)
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(ValidNetworkProtocols, stringValue) {
		return networkProtocol, errors.New("InvalidNetworkProtocol")
	}

	return NetworkProtocol(stringValue), nil
}

func (vo NetworkProtocol) String() string {
	return string(vo)
}
