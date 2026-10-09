package infraHelper

import (
	"errors"
	"log/slog"
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
	for attemptNumber := range nAttempts {
		_, err := tkInfra.NewShell(tkInfra.ShellSettings{
			Command: "apt-get",
			Args:    installPackages,
		}).Run()
		if err == nil {
			break
		}

		slog.Debug("InstallPkgError", slog.String("err", err.Error()))

		if attemptNumber == nAttempts-1 {
			installErr = errors.New("InstallAttemptsFailed")
		}
	}

	for _, aptCacheDir := range []string{
		"/var/lib/apt/lists", "/var/cache/apt/archives",
	} {
		err := os.RemoveAll(aptCacheDir)
		if err != nil {
			slog.Debug(
				"AptCacheCleanupFailed",
				slog.String("aptCacheDir", aptCacheDir),
				slog.String("err", err.Error()),
			)
		}
	}

	return installErr
}
