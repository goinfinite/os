package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ServiceStatus string

var ServiceStatusesWithAliases = map[string][]string{
	"running": {
		"run", "up", "start", "started", "enable", "enabled", "activate", "active",
		"true", "on", "ok", "yes", "y", "1",
	},
	"stopped": {
		"stop", "halt", "halted", "pause", "paused", "deactivate", "deactivated",
		"false", "off", "no", "n", "0",
	},
	"uninstalled": {
		"uninstall", "uninstalled", "remove", "removed", "delete", "deleted",
		"purge", "purged", "clear", "cleared", "clean", "cleaned",
	},
	"restarting": {
		"restart", "restarted", "reload", "reloaded", "refresh", "refreshed",
		"reboot", "rebooted", "reset", "reseted",
	},
}

func NewServiceStatus(value interface{}) (status ServiceStatus, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return status, errors.New("ServiceStatusMustBeString")
	}

	stringValue, err = serviceStatusAdapter(stringValue)
	if err != nil {
		return status, err
	}

	return ServiceStatus(stringValue), nil
}

func serviceStatusAdapter(value string) (string, error) {
	value = strings.ToLower(value)

	for exactName, aliases := range ServiceStatusesWithAliases {
		if exactName == value {
			return exactName, nil
		}

		if slices.Contains(aliases, value) {
			return exactName, nil
		}
	}

	return "", errors.New("InvalidServiceStatus")
}

func (vo ServiceStatus) String() string {
	return string(vo)
}
