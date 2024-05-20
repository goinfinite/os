package infraHelper

import "errors"

func ValidateWebServerConfig() error {
	_, err := RunCmdWithSubShell("nginx -t")
	return errors.New("WebServerConfigValidationError: " + err.Error())
}
