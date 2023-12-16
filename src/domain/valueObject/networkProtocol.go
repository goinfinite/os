package valueObject

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type NetworkProtocol string

var ValidNetworkProtocols = []string{
	"http",
	"https",
	"unix",
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
