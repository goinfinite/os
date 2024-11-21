package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadInstallableServicesItemsRequest struct {
	Pagination    Pagination                 `json:"pagination"`
	ServiceName   *valueObject.ServiceName   `json:"name,omitempty"`
	ServiceNature *valueObject.ServiceNature `json:"nature,omitempty"`
	ServiceType   *valueObject.ServiceType   `json:"type,omitempty"`
}

type ReadInstallableServicesItemsResponse struct {
	Pagination          Pagination                  `json:"pagination"`
	InstallableServices []entity.InstallableService `json:"installableServices"`
}
