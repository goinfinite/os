package infraHelper

import "github.com/goinfinite/os/src/domain/valueObject"

func IsPrimaryVirtualHost(vhost valueObject.Fqdn) bool {
	primaryVhost, err := GetPrimaryVirtualHost()
	if err != nil {
		return false
	}

	return vhost == primaryVhost
}
