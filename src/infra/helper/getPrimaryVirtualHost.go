package infraHelper

import (
	"os"

	tkInfra "github.com/goinfinite/tk/src/infra"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func ReadPrimaryVirtualHostHostname() (primaryHostname tkValueObject.Fqdn, err error) {
	rawPrimaryHost := os.Getenv("PRIMARY_VHOST")
	if rawPrimaryHost != "" {
		return tkValueObject.NewFqdn(rawPrimaryHost)
	}

	rawPrimaryHost, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "hostname",
		Args:    []string{"-f"},
	}).Run()
	if err != nil {
		return primaryHostname, err
	}

	return tkValueObject.NewFqdn(rawPrimaryHost)
}
