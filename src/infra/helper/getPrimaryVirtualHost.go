package infraHelper

import (
	"os"

	"github.com/speedianet/os/src/domain/valueObject"
)

func GetPrimaryVirtualHost() (valueObject.Fqdn, error) {
	var primaryVhost valueObject.Fqdn

	primaryVhostStr := os.Getenv("PRIMARY_VHOST")
	if primaryVhostStr == "" {
		var err error
		primaryVhostStr, err = RunCmd("hostname", "-f")
		if err != nil {
			return primaryVhost, err
		}
	}

	return valueObject.NewFqdn(primaryVhostStr)
}
