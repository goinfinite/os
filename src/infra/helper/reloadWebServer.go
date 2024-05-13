package infraHelper

import "errors"

func ReloadWebServer() error {
	_, err := RunCmdWithSubShell(
		"nginx -t && nginx -s reload && sleep 2",
	)
	if err != nil {
		return errors.New("NginxReloadFailed: " + err.Error())
	}

	return nil
}
