package valueObject

import (
	"errors"
	"net"
	"regexp"
)

const fqdnRegex string = `^((\*\.)?([a-zA-Z0-9_]+[\w-]*\.)*)?([a-zA-Z0-9_][\w-]*[a-zA-Z0-9])\.([a-zA-Z]{2,})$`

type Fqdn string

func NewFqdn(value string) (Fqdn, error) {
	fqdn := Fqdn(value)
	if !fqdn.isValid() {
		return "", errors.New("InvalidFqdn")
	}
	return fqdn, nil
}

func NewFqdnPanic(value string) Fqdn {
	fqdn := Fqdn(value)
	if !fqdn.isValid() {
		panic("InvalidFqdn")
	}
	return fqdn
}

func (fqdn Fqdn) isValid() bool {
	if net.ParseIP(string(fqdn)) != nil {
		return false
	}
	re := regexp.MustCompile(fqdnRegex)
	return re.MatchString(string(fqdn))
}

func (fqdn Fqdn) String() string {
	return string(fqdn)
}
