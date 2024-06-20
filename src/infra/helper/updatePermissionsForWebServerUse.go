package infraHelper

import "strings"

func UpdatePermissionsForWebServerUse(
	filePath string, isRecursive bool, shouldIncludeSymlink bool,
) error {
	flags := []string{}
	if isRecursive {
		flags = append(flags, "-R")
	}

	if shouldIncludeSymlink {
		flags = append(flags, "-L")
	}
	flagsStr := strings.Join(flags, " ")

	params := []string{}
	webServerUsername := "nobody"
	webServerUserGroup := "nogroup"
	params = append(params, webServerUsername+":"+webServerUserGroup)

	params = append(params, filePath)

	paramsStr := strings.Join(params, " ")
	_, err := RunCmdWithSubShell("chown " + flagsStr + " " + paramsStr)
	if err != nil {
		return err
	}

	return nil
}
