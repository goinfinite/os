package infraHelper

import "errors"

func ValidateWebServerConfig() error {
	_, err := RunCmd(RunCmdConfigs{
		Command:               "nginx -t",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("WebServerConfigValidationError: " + err.Error())
	}

	return nil
}
