package infraHelper

import (
	"errors"

	"github.com/goinfinite/os/src/domain/valueObject"
	"golang.org/x/net/publicsuffix"
)

func GetRootDomain(hostname valueObject.Fqdn) (valueObject.Fqdn, error) {
	var rootDomain valueObject.Fqdn

	rootDomainStr, err := publicsuffix.EffectiveTLDPlusOne(hostname.String())
	if err != nil {
		return rootDomain, errors.New("InvalidHostname")
	}

	return valueObject.NewFqdn(rootDomainStr)
}
