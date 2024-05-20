package infraHelper

import "errors"

func ValidateWebServerConfig() error {
	_, err := RunCmdWithSubShell("nginx -t")
	if err != nil {
		return errors.New("WebServerConfigValidationError: " + err.Error())
	}

	return nil
}
