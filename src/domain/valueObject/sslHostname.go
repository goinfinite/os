package valueObject

import (
	"errors"
	"net"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

type SslHostname string

func NewSslHostname(value interface{}) (sslHostname SslHostname, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return sslHostname, errors.New("SslHostnameMustBeString")
	}

	_, err = NewFqdn(stringValue)
	if err == nil {
		return SslHostname(stringValue), nil
	}

	ipAddress := net.ParseIP(stringValue)
	if ipAddress != nil {
		return SslHostname(stringValue), nil
	}

	return sslHostname, errors.New("InvalidSslHostname")
}

func (vo SslHostname) String() string {
	return string(vo)
}
