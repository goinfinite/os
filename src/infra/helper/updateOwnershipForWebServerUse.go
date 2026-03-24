package infraHelper

import (
	"strings"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func UpdateOwnershipForWebServerUse(
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
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           "chown " + flagsStr + " " + paramsStr,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return err
	}

	return nil
}
