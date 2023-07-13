package dto

import "github.com/speedianet/sam/src/domain/valueObject"

type UpdateSvcStatus struct {
	Name    valueObject.ServiceName     `json:"name"`
	Status  valueObject.ServiceStatus   `json:"status"`
	Version *valueObject.ServiceVersion `json:"version,omitempty"`
}

func NewUpdateSvcStatus(
	name valueObject.ServiceName,
	status valueObject.ServiceStatus,
	version *valueObject.ServiceVersion,
) UpdateSvcStatus {
	return UpdateSvcStatus{
		Name:    name,
		Status:  status,
		Version: version,
	}
}
