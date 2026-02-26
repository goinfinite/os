package infraHelper

import (
	"strings"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

type RunCmdSettings struct {
	Command               string
	Args                  []string
	Username              string
	WorkingDirectory      string
	ShouldRunWithSubShell bool
	ExecutionTimeoutSecs  uint64
}

func IsRunCmdTimeout(err error) bool {
	return strings.Contains(err.Error(), "CommandDeadlineExceeded")
}

func RunCmd(runSettings RunCmdSettings) (string, error) {
	return tkInfra.NewShell(tkInfra.ShellSettings{
		Command:              runSettings.Command,
		Args:                 runSettings.Args,
		Username:             runSettings.Username,
		WorkingDirectory:     runSettings.WorkingDirectory,
		ShouldUseSubShell:    runSettings.ShouldRunWithSubShell,
		ExecutionTimeoutSecs: runSettings.ExecutionTimeoutSecs,
	}).Run()
}
