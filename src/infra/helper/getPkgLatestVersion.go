package infraHelper

import (
	"errors"
	"fmt"
	"strings"
)

func GetPkgLatestVersion(pkgName string, majorVersion *string) (string, error) {
	out, err := RunCmd("apt", "list", "-a", pkgName)
	if err != nil {
		return "", err
	}

	lines := strings.Split(out, "\n")
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
		if majorVersion == nil {
			return currentVersion, nil
		}

		if strings.HasPrefix(line, fmt.Sprintf("%s/%s", pkgName, *majorVersion)) {
			return currentVersion, nil
		}
	}

	return "", errors.New("PackageOrVersionNotFound")
}
