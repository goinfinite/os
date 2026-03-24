package infraHelper

import (
	"errors"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func ValidateWebServerConfig() error {
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:          "nginx -t",
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return errors.New("WebServerConfigValidationError: " + err.Error())
	}

	return nil
}
