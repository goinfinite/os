package valueObject

import (
	"errors"
	"slices"
	"strings"
)

type NetworkProtocol string

var ValidNetworkProtocols = []string{
	"http",
	"https",
	"ws",
	"wss",
	"grpc",
	"grpcs",
	"tcp",
	"udp",
}

func NewNetworkProtocol(value string) (NetworkProtocol, error) {
	value = strings.ToLower(value)
	if !slices.Contains(ValidNetworkProtocols, value) {
		return "", errors.New("InvalidNetworkProtocol")
	}
	return NetworkProtocol(value), nil
}

func NewNetworkProtocolPanic(value string) NetworkProtocol {
	np, err := NewNetworkProtocol(value)
	if err != nil {
		panic(err)
	}
	return np
}

func (np NetworkProtocol) String() string {
	return string(np)
}
