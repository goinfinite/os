package infraHelper

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

const commandDeadlineExceededError string = "CommandDeadlineExceeded"

func IsRunCmdTimeout(err error) bool {
	return strings.Contains(err.Error(), commandDeadlineExceededError)
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
	ExecutionTimeoutSecs  uint64
}

func prepareCmdExecutor(
	runConfigs RunCmdConfigs,
) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	args := runConfigs.Args
	command := runConfigs.Command
	if runConfigs.ShouldRunWithSubShell {
		args = []string{
			"-c", "source /etc/profile; " + command + " " + strings.Join(args, " "),
		}
		command = "bash"
	}

	executionTimeoutSecs := uint64(1800)
	if runConfigs.ExecutionTimeoutSecs > 0 && runConfigs.ExecutionTimeoutSecs <= executionTimeoutSecs {
		executionTimeoutSecs = runConfigs.ExecutionTimeoutSecs
	}

	args = slices.Concat(
		[]string{strconv.FormatUint(runConfigs.ExecutionTimeoutSecs, 10), command},
		args,
	)
	command = "timeout"

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
				stdErrStr = commandDeadlineExceededError
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
