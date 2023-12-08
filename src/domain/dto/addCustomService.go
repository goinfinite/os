package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddCustomService struct {
	Name    valueObject.ServiceName     `json:"name"`
	Type    valueObject.ServiceType     `json:"type"`
	Command valueObject.UnixCommand     `json:"command"`
	Version *valueObject.ServiceVersion `json:"version,omitempty"`
	Ports   []valueObject.NetworkPort   `json:"ports,omitempty"`
}

func NewAddCustomService(
	name valueObject.ServiceName,
	serviceType valueObject.ServiceType,
	command valueObject.UnixCommand,
	version *valueObject.ServiceVersion,
	ports []valueObject.NetworkPort,
) AddCustomService {
	return AddCustomService{
		Name:    name,
		Type:    serviceType,
		Command: command,
		Version: version,
		Ports:   ports,
	}
}
