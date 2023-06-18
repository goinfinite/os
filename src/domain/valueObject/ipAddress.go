package valueObject

import (
	"errors"
	"net"
)

type IpAddress string

func NewIpAddress(value string) (IpAddress, error) {
	addr := IpAddress(value)
	if !addr.isValid(value) {
		return "", errors.New("InvalidIpAddress")
	}
	return addr, nil
}

func NewIpAddressPanic(value string) IpAddress {
	addr := IpAddress(value)
	if !addr.isValid(value) {
		panic("InvalidIpAddress")
	}
	return addr
}

func (addr IpAddress) isValid(value string) bool {
	ip := net.ParseIP(value)
	return ip != nil
}

func (addr IpAddress) String() string {
	return string(addr)
}
