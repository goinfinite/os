package infraHelper

import "errors"

func ValidateWebServerConfig() error {
	_, err := RunCmd(RunCmdSettings{
		Command:               "nginx -t",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("WebServerConfigValidationError: " + err.Error())
	}

	return nil
}
