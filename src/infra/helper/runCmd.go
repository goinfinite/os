package infraHelper

import (
	"bytes"
	"encoding/json"
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

func RunCmd(command string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmdObj := exec.Command(command, args...)
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
