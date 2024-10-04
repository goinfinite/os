package infraHelper

import (
	"errors"
	"strings"
	"time"
)

func ReloadWebServer() error {
	wsConfigTestResult, err := RunCmd("/usr/sbin/nginx", "-t")
	if err != nil {
		if wsConfigTestResult != "" {
			return errors.New("NginxConfigTestFailed: " + err.Error())
		}

		if !strings.Contains(wsConfigTestResult, "test is successful") {
			return errors.New("NginxConfigTestFailed: " + wsConfigTestResult)
		}
	}

	_, err = RunCmd("/usr/sbin/nginx", "-s", "reload", "-c", "/etc/nginx/nginx.conf")
	if err != nil {
		return errors.New("NginxReloadFailed: " + err.Error())
	}

	time.Sleep(2 * time.Second)

	return nil
}
