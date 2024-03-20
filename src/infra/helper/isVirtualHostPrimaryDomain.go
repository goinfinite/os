package infraHelper

import "github.com/speedianet/os/src/domain/valueObject"

func IsVirtualHostPrimaryDomain(domain valueObject.Fqdn) bool {
	primaryHostname, err := GetPrimaryHostname()
	if err != nil {
		return false
	}

	return domain == primaryHostname
}
