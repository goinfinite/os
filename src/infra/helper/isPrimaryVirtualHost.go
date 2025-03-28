package infraHelper

import "github.com/goinfinite/os/src/domain/valueObject"

func IsPrimaryVirtualHost(vhost valueObject.Fqdn) bool {
	primaryVhost, err := ReadPrimaryVirtualHostHostname()
	if err != nil {
		return false
	}

	return vhost == primaryVhost
}
