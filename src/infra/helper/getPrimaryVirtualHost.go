package infraHelper

import (
	"os"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func ReadPrimaryVirtualHostHostname() (primaryHostname tkValueObject.Fqdn, err error) {
	rawPrimaryHost := os.Getenv("PRIMARY_VHOST")
	if rawPrimaryHost != "" {
		return tkValueObject.NewFqdn(rawPrimaryHost)
	}

	rawPrimaryHost, err = RunCmd(RunCmdSettings{
		Command: "hostname",
		Args:    []string{"-f"},
	})
	if err != nil {
		return primaryHostname, err
	}

	return tkValueObject.NewFqdn(rawPrimaryHost)
}
