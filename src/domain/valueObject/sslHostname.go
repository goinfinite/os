package valueObject

import (
	"errors"
	"net"
)

type SslHostname string

func NewSslHostname(value string) (SslHostname, error) {
	sslHostname := SslHostname(value)
	if !sslHostname.isValid() {
		return "", errors.New("InvalidSslHostname")
	}

	return sslHostname, nil
}

func NewSslHostnamePanic(value string) SslHostname {
	sslHostname, err := NewSslHostname(value)
	if err != nil {
		panic(err)
	}

	return sslHostname
}

func (sslHostname SslHostname) isValid() bool {
	_, err := NewFqdn(string(sslHostname))
	if err == nil {
		return true
	}

	return net.ParseIP(string(sslHostname)) != nil
}

func (sslHostname SslHostname) String() string {
	return string(sslHostname)
}
