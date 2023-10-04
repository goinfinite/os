package valueObject

import (
	"errors"
	"regexp"
)

const virtualHostRegex string = `^[a-z]{1}[0-9a-zA-Z_\.-]{2,63}$`

type VirtualHost string

func NewVirtualHost(value string) (VirtualHost, error) {
	virtualHost := VirtualHost(value)
	if !virtualHost.isValid() {
		return "", errors.New("InvalidVirtualHost")
	}
	return virtualHost, nil
}

func NewVirtualHostPanic(value string) VirtualHost {
	virtualHost := VirtualHost(value)
	if !virtualHost.isValid() {
		panic("InvalidVirtualHost")
	}
	return virtualHost
}

func (virtualHost VirtualHost) isValid() bool {
	re := regexp.MustCompile(virtualHostRegex)
	return re.MatchString(string(virtualHost))
}

func (virtualHost VirtualHost) String() string {
	return string(virtualHost)
}
