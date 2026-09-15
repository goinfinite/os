package infraHelper

import (
	"errors"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func RemovePkgs(packages []string) error {
	removePackages := append([]string{"remove", "-y"}, packages...)

	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "apt-get",
		Args:    removePackages,
	}).Run()
	if err != nil {
		return errors.New("RemovePkgsFailed: " + err.Error())
	}

	return nil
}
