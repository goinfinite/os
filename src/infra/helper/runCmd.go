package infraHelper

import (
	"bytes"
	"encoding/json"
	"os/exec"
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
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return stdout.String(), &CommandError{
				StdErr:   stderr.String(),
				ExitCode: exitErr.ExitCode(),
			}
		}
		return stdout.String(), err
	}

	return stdout.String(), nil
}
