package infraHelper

import (
	"errors"
	"log"
	"os"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func InstallPkgs(packages []string) error {
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "apt-get",
		Args:    []string{"update", "-qq"},
	}).Run()
	if err != nil {
		return errors.New("UpdateRepositoriesFailed")
	}

	installPackages := append(
		[]string{"install", "-y", "--no-install-recommends"},
		packages...,
	)

	var installErr error
	nAttempts := 3
	for i := 0; i < nAttempts; i++ {
		_, err := tkInfra.NewShell(tkInfra.ShellSettings{
			Command: "apt-get",
			Args:    installPackages,
		}).Run()
		if err == nil {
			break
		}

		log.Printf("InstallPkgError: %s", err.Error())

		if i == nAttempts-1 {
			installErr = errors.New("InstallAttemptsFailed")
		}
	}

	os.RemoveAll("/var/lib/apt/lists")
	os.RemoveAll("/var/cache/apt/archives")

	return installErr
}
