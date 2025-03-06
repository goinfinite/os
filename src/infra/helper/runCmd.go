package infraHelper

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strings"
)

const CommandDeadlineExceededError string = "CommandDeadlineExceeded"

func IsRunCmdTimeout(err error) bool {
	return strings.Contains(err.Error(), CommandDeadlineExceededError)
}

type CmdError struct {
	StdErr   string `json:"stdErr"`
	ExitCode int    `json:"exitCode"`
}

func (e *CmdError) Error() string {
	jsonError, _ := json.Marshal(e)
	return string(jsonError)
}

type RunCmdConfigs struct {
	Command               string
	Args                  []string
	ShouldRunWithSubShell bool
	CmdContext            context.Context
}

func prepareCmdExecutor(
	runConfigs RunCmdConfigs,
) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	command := runConfigs.Command
	args := runConfigs.Args
	if runConfigs.ShouldRunWithSubShell {
		subShellCommand := "bash"
		subShellArgs := []string{"-c", "source /etc/profile; " + command}

		command = subShellCommand
		args = subShellArgs
	}

	cmdExecutor := exec.Command(command, args...)
	if runConfigs.CmdContext != nil {
		cmdExecutor = exec.CommandContext(runConfigs.CmdContext, command, args...)
	}

	var stdoutBytesBuffer, stderrBytesBuffer bytes.Buffer
	cmdExecutor.Stdout = &stdoutBytesBuffer
	cmdExecutor.Stderr = &stderrBytesBuffer

	cmdExecutor.Env = append(cmdExecutor.Environ(), "DEBIAN_FRONTEND=noninteractive")

	return cmdExecutor, &stdoutBytesBuffer, &stderrBytesBuffer
}

func RunCmd(runConfigs RunCmdConfigs) (string, error) {
	cmdExecutor, stdoutBytesBuffer, stderrBytesBuffer := prepareCmdExecutor(runConfigs)

	err := cmdExecutor.Run()
	stdoutStr := strings.TrimSpace(stdoutBytesBuffer.String())

	if runConfigs.CmdContext != nil {
		if runConfigs.CmdContext.Err() == context.DeadlineExceeded {
			return stdoutStr, &CmdError{
				StdErr:   CommandDeadlineExceededError,
				ExitCode: 124,
			}
		}
	}

	if err != nil {
		if exitErr, assertOk := err.(*exec.ExitError); assertOk {
			return stdoutStr, &CmdError{
				StdErr:   stderrBytesBuffer.String(),
				ExitCode: exitErr.ExitCode(),
			}
		}
		return stdoutStr, err
	}

	return stdoutStr, nil
}
