package infraHelper

import (
	"errors"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	"golang.org/x/net/publicsuffix"
)

func GetRootDomain(hostname tkValueObject.Fqdn) (tkValueObject.Fqdn, error) {
	var rootDomain tkValueObject.Fqdn

	rootDomainStr, err := publicsuffix.EffectiveTLDPlusOne(hostname.String())
	if err != nil {
		return rootDomain, errors.New("InvalidHostname")
	}

	return tkValueObject.NewFqdn(rootDomainStr)
}
