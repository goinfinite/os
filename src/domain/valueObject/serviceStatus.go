package valueObject

import (
	"errors"
	"slices"
	"strings"

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
}

func NewServiceStatus(value string) (ServiceStatus, error) {
	value, err := ServiceStatusAdapter(value)
	if err != nil {
		return "", err
	}

	return ServiceStatus(value), nil
}

func NewServiceStatusPanic(value string) ServiceStatus {
	ss, err := NewServiceStatus(value)
	if err != nil {
		panic(err)
	}
	return ss
}

func ServiceStatusAdapter(value string) (string, error) {
	value = strings.TrimSpace(value)
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
		return "", errors.New("InvalidServiceStatus")
	}

	return value, nil
}

func (vo ServiceStatus) String() string {
	return string(vo)
}
