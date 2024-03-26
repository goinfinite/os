package infraHelper

import "github.com/speedianet/os/src/domain/valueObject"

func IsPrimaryVirtualHost(host valueObject.Fqdn) bool {
	primaryHost, err := GetPrimaryHostname()
	if err != nil {
		return false
	}

	return host == primaryHost
}
