package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

type ReadInstallableServicesItemsRequest struct {
	Pagination    tkDto.Pagination           `json:"pagination"`
	ServiceName   *valueObject.ServiceName   `json:"name,omitempty"`
	ServiceNature *valueObject.ServiceNature `json:"nature,omitempty"`
	ServiceType   *valueObject.ServiceType   `json:"type,omitempty"`
}

type ReadInstallableServicesItemsResponse struct {
	Pagination          tkDto.Pagination            `json:"pagination"`
	InstallableServices []entity.InstallableService `json:"installableServices"`
}
