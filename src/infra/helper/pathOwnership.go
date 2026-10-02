package infraHelper

import (
	"strings"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type PathOwnership struct{}

func (ownership PathOwnership) update(
	filePath, ownerUsername, ownerGroupName string,
	isRecursive, shouldIncludeSymlink bool,
) error {
	flags := []string{}
	if isRecursive {
		flags = append(flags, "-R")
	}

	if shouldIncludeSymlink {
		flags = append(flags, "-L")
	}
	flagsStr := strings.Join(flags, " ")

	paramsStr := strings.Join(
		[]string{ownerUsername + ":" + ownerGroupName, filePath}, " ",
	)

	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           "chown " + flagsStr + " " + paramsStr,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return err
	}

	return nil
}

func (ownership PathOwnership) UpdateForWebServerUse(
	filePath string, isRecursive bool, shouldIncludeSymlink bool,
) error {
	return ownership.update(
		filePath, infraEnvs.PhpWebServerUsername, infraEnvs.PhpWebServerGroupName,
		isRecursive, shouldIncludeSymlink,
	)
}

func (ownership PathOwnership) UpdateForRootUse(
	filePath, groupName string, isRecursive bool,
) error {
	return ownership.update(
		filePath, "root", groupName, isRecursive, false,
	)
}
