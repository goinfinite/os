package valueObject

import (
	"errors"
	"net"
	"regexp"
	"strings"

	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
)

const fqdnRegex string = `^((\*\.)?([a-zA-Z0-9_]+[\w-]*\.)*)?([a-zA-Z0-9_][\w-]*[a-zA-Z0-9])\.([a-zA-Z]{2,})$`

type Fqdn string

func NewFqdn(value interface{}) (fqdn Fqdn, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return fqdn, errors.New("FqdnMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	isIpAddress := net.ParseIP(stringValue) != nil
	if isIpAddress {
		return fqdn, errors.New("FqdnCannotBeIpAddress")
	}

	re := regexp.MustCompile(fqdnRegex)
	if !re.MatchString(stringValue) {
		return fqdn, errors.New("InvalidFqdn")
	}

	return Fqdn(stringValue), nil
}

func (vo Fqdn) String() string {
	return string(vo)
}
