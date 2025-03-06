package infraHelper

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"slices"
	"strconv"
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
	ExecutionTimeout      uint64
}

func prepareCmdExecutor(
	runConfigs RunCmdConfigs,
) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	command := runConfigs.Command
	args := runConfigs.Args
	if runConfigs.ShouldRunWithSubShell {
		subShellCommand := "bash"

		argsStr := strings.Join(args, " ")
		subShellArgs := []string{
			"-c", "source /etc/profile; " + command + " " + argsStr,
		}

		command = subShellCommand
		args = subShellArgs
	}

	if runConfigs.ExecutionTimeout > 0 {
		timeoutCommand := "timeout"
		timeoutArgs := []string{
			strconv.FormatUint(runConfigs.ExecutionTimeout, 10), command,
		}

		command = timeoutCommand
		args = slices.Concat(timeoutArgs, args)
	}

	cmdExecutor := exec.Command(command, args...)

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

	if err != nil {
		if exitErr, assertOk := err.(*exec.ExitError); assertOk {
			stdErrStr := stderrBytesBuffer.String()
			if exitErr.ExitCode() == 124 {
				stdErrStr = CommandDeadlineExceededError
			}

			return stdoutStr, &CmdError{
				StdErr:   stdErrStr,
				ExitCode: exitErr.ExitCode(),
			}
		}
		return stdoutStr, err
	}

	return stdoutStr, nil
}
