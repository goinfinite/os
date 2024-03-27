package infraHelper

import "github.com/speedianet/os/src/domain/valueObject"

func IsPrimaryVirtualHost(vhost valueObject.Fqdn) bool {
	primaryVhost, err := GetPrimaryVirtualHost()
	if err != nil {
		return false
	}

	return vhost == primaryVhost
}
