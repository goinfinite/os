package infraHelper

import (
	"errors"

	"github.com/speedianet/os/src/domain/valueObject"
	"golang.org/x/net/publicsuffix"
)

func GetRootDomain(serverName valueObject.Fqdn) (valueObject.Fqdn, error) {
	var rootDomain valueObject.Fqdn

	rootDomainStr, err := publicsuffix.EffectiveTLDPlusOne(serverName.String())
	if err != nil {
		return rootDomain, errors.New("InvalidRootDomain")
	}

	rootDomain, err = valueObject.NewFqdn(rootDomainStr)
	if err != nil {
		return rootDomain, errors.New("InvalidRootDomain")
	}

	return rootDomain, nil
}
