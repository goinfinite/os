package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadInstallableServicesItemsRequest struct {
	Pagination Pagination                 `json:"pagination"`
	Name       *valueObject.ServiceName   `json:"name,omitempty"`
	Nature     *valueObject.ServiceNature `json:"nature,omitempty"`
	Type       *valueObject.ServiceType   `json:"type,omitempty"`
}

type ReadInstallableServicesItemsResponse struct {
	Pagination Pagination                  `json:"pagination"`
	Items      []entity.InstallableService `json:"items"`
}
