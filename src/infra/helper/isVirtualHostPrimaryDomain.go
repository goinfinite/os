package infraHelper

import "github.com/speedianet/os/src/domain/valueObject"

func IsVirtualHostPrimaryDomain(domain valueObject.Fqdn) bool {
	primaryDomain, err := GetPrimaryHostname()
	if err != nil {
		return false
	}

	return domain == primaryDomain
}
