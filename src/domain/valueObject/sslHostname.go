package valueObject

import (
	"errors"
	"net"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type SslHostname string

func NewSslHostname(value interface{}) (sslHostname SslHostname, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return sslHostname, errors.New("SslHostnameMustBeString")
	}

	_, err = tkValueObject.NewFqdn(stringValue)
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
