package infraHelper

import (
	"os"

	"github.com/speedianet/os/src/domain/valueObject"
)

func GetPrimaryHostname() (valueObject.Fqdn, error) {
	var hostname valueObject.Fqdn

	hostnameStr := os.Getenv("VIRTUAL_HOST")
	if hostnameStr == "" {
		var err error
		hostnameStr, err = RunCmd("hostname", "-f")
		if err != nil {
			return hostname, err
		}
	}

	return valueObject.NewFqdn(hostnameStr)
}
