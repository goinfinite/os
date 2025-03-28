package infraHelper

import (
	"os"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func ReadPrimaryVirtualHostHostname() (primaryHostname valueObject.Fqdn, err error) {
	rawPrimaryHost := os.Getenv("PRIMARY_VHOST")
	if rawPrimaryHost != "" {
		return valueObject.NewFqdn(rawPrimaryHost)
	}

	rawPrimaryHost, err = RunCmd(RunCmdSettings{
		Command: "hostname",
		Args:    []string{"-f"},
	})
	if err != nil {
		return primaryHostname, err
	}

	return valueObject.NewFqdn(rawPrimaryHost)
}
