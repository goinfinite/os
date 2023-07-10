package valueObject

import "errors"

type ServiceStatus string

const (
	running     ServiceStatus = "running"
	stopped     ServiceStatus = "stopped"
	uninstalled ServiceStatus = "uninstalled"
	installing  ServiceStatus = "installing"
)

func NewServiceStatus(value string) (ServiceStatus, error) {
	ss := ServiceStatus(value)
	if !ss.isValid() {
		return "", errors.New("InvalidServiceStatus")
	}
	return ss, nil
}

func NewServiceStatusPanic(value string) ServiceStatus {
	ss := ServiceStatus(value)
	if !ss.isValid() {
		panic("InvalidServiceStatus")
	}
	return ss
}

func (ss ServiceStatus) isValid() bool {
	switch ss {
	case running, stopped, uninstalled, installing:
		return true
	default:
		return false
	}
}

func (ss ServiceStatus) String() string {
	return string(ss)
}
