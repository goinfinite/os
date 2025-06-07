package infraHelper

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"os/user"
	"slices"
	"strconv"
	"strings"
	"syscall"
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

type RunCmdSettings struct {
	Command               string
	Args                  []string
	Username              string
	WorkingDirectory      string
	ShouldRunWithSubShell bool
	ExecutionTimeoutSecs  uint64
}

func sysCallCredentialsFactory(
	username string,
) (*syscall.Credential, error) {
	userStruct, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}
	userId, err := strconv.Atoi(userStruct.Uid)
	if err != nil {
		return nil, err
	}
	groupId, err := strconv.Atoi(userStruct.Gid)
	if err != nil {
		return nil, err
	}

	return &syscall.Credential{
		Uid: uint32(userId),
		Gid: uint32(groupId),
	}, nil
}

func prepareCmdExecutor(
	runSettings RunCmdSettings,
) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	args := runSettings.Args
	command := runSettings.Command
	if runSettings.ShouldRunWithSubShell {
		args = []string{
			"-c", "source /etc/profile; " + command + " " + strings.Join(args, " "),
		}
		command = "bash"
	}

	executionTimeoutSecs := uint64(1800)
	if runSettings.ExecutionTimeoutSecs > 0 && runSettings.ExecutionTimeoutSecs <= executionTimeoutSecs {
		executionTimeoutSecs = runSettings.ExecutionTimeoutSecs
	}

	args = slices.Concat(
		[]string{strconv.FormatUint(executionTimeoutSecs, 10), command},
		args,
	)
	command = "timeout"

	cmdExecutor := exec.Command(command, args...)
	if runSettings.Username != "" && runSettings.Username != "root" {
		sysCallCredentials, err := sysCallCredentialsFactory(runSettings.Username)
		if err == nil {
			cmdExecutor.SysProcAttr = &syscall.SysProcAttr{Credential: sysCallCredentials}
		}
	}

	if runSettings.WorkingDirectory != "" {
		cmdExecutor.Dir = runSettings.WorkingDirectory
	}

	var stdoutBytesBuffer, stderrBytesBuffer bytes.Buffer
	cmdExecutor.Stdout = &stdoutBytesBuffer
	cmdExecutor.Stderr = &stderrBytesBuffer

	cmdExecutor.Env = append(cmdExecutor.Environ(), "DEBIAN_FRONTEND=noninteractive")

	return cmdExecutor, &stdoutBytesBuffer, &stderrBytesBuffer
}

func RunCmd(runSettings RunCmdSettings) (string, error) {
	cmdExecutor, stdoutBytesBuffer, stderrBytesBuffer := prepareCmdExecutor(runSettings)

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
