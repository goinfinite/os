package infraHelper

import (
	"errors"
	"os"
	"strings"
)

func GetPkgLatestVersion(pkgName string, majorVersion *string) (string, error) {
	_, _ = RunCmd("apt", "update")
	out, err := RunCmd("apt", "list", "-a", pkgName)
	if err != nil {
		return "", err
	}
	os.RemoveAll("/var/lib/apt/lists")
	os.RemoveAll("/var/cache/apt/archives")

	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		return "", errors.New("PackageNotFound")
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, pkgName) {
			continue
		}

		versionDetails := strings.Fields(line)
		if len(versionDetails) < 2 {
			continue
		}
		currentVersion := strings.Fields(line)[1]
		if majorVersion == nil || strings.Contains(line, *majorVersion) {
			return currentVersion, nil
		}
	}

	return "", errors.New("VersionNotFound")
}
