package infraHelper

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

type CommandError struct {
	StdErr   string `json:"stdErr"`
	ExitCode int    `json:"exitCode"`
}

func (e *CommandError) Error() string {
	errJSON, _ := json.Marshal(e)
	return string(errJSON)
}

func execCmd(cmdObj *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmdObj.Stdout = &stdout
	cmdObj.Stderr = &stderr

	err := cmdObj.Run()
	stdOut := strings.TrimSpace(stdout.String())
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return stdOut, &CommandError{
				StdErr:   stderr.String(),
				ExitCode: exitErr.ExitCode(),
			}
		}
		return stdOut, err
	}

	return stdOut, nil
}

func RunCmd(command string, args ...string) (string, error) {
	cmdObj := exec.Command(command, args...)
	return execCmd(cmdObj)
}

func RunCmdWithEnvVars(
	command string,
	envVars map[string]string,
	args ...string,
) (string, error) {
	cmdObj := exec.Command(command, args...)
	cmdObj.Env = os.Environ()
	for envVar, envValue := range envVars {
		cmdObj.Env = append(cmdObj.Env, envVar+"="+envValue)
	}
	cmdObj.Env = append(cmdObj.Env)
	return execCmd(cmdObj)
}
