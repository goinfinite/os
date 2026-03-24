package infraHelper

import (
	"errors"
	"strings"
	"time"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func ReloadWebServer() error {
	wsConfigTestResult, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "/usr/sbin/nginx",
		Args:    []string{"-t"},
	}).Run()
	if err != nil {
		if wsConfigTestResult != "" {
			return errors.New("NginxConfigTestFailed: " + err.Error())
		}

		if !strings.Contains(wsConfigTestResult, "test is successful") {
			return errors.New("NginxConfigTestFailed: " + wsConfigTestResult)
		}
	}

	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "/usr/sbin/nginx",
		Args:    []string{"-s", "reload", "-c", "/etc/nginx/nginx.conf"},
	}).Run()
	if err != nil {
		return errors.New("NginxReloadFailed: " + err.Error())
	}

	time.Sleep(2 * time.Second)

	return nil
}
