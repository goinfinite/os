package infraHelper

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

func IsPrimaryVirtualHost(vhost tkValueObject.Fqdn) bool {
	primaryVhost, err := ReadPrimaryVirtualHostHostname()
	if err != nil {
		return false
	}

	return vhost == primaryVhost
}
