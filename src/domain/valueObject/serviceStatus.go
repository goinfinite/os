package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	"golang.org/x/exp/maps"
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

	if _, isPrimaryStatus := ServiceStatusesWithAliases[value]; isPrimaryStatus {
		return value, nil
	}

	primaryStatuses := maps.Keys(ServiceStatusesWithAliases)
	for _, primaryStatus := range primaryStatuses {
		if !slices.Contains(ServiceStatusesWithAliases[primaryStatus], value) {
			continue
		}
		value = primaryStatus
		break
	}

	if _, isPrimaryStatus := ServiceStatusesWithAliases[value]; !isPrimaryStatus {
		return value, errors.New("InvalidServiceStatus")
	}

	return value, nil
}

func (vo ServiceStatus) String() string {
	return string(vo)
}
