package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

type ReadInstalledServicesItemsRequest struct {
	Pagination           tkDto.Pagination           `json:"pagination"`
	ServiceName          *valueObject.ServiceName   `json:"serviceName,omitempty"`
	ServiceNature        *valueObject.ServiceNature `json:"serviceNature,omitempty"`
	ServiceStatus        *valueObject.ServiceStatus `json:"serviceStatus,omitempty"`
	ServiceType          *valueObject.ServiceType   `json:"serviceType,omitempty"`
	ShouldIncludeMetrics *bool                      `json:"shouldIncludeMetrics,omitempty"`
}

type ReadFirstInstalledServiceItemsRequest struct {
	Pagination    tkDto.Pagination           `json:"pagination"`
	ServiceName   *valueObject.ServiceName   `json:"serviceName,omitempty"`
	ServiceNature *valueObject.ServiceNature `json:"serviceNature,omitempty"`
	ServiceStatus *valueObject.ServiceStatus `json:"serviceStatus,omitempty"`
	ServiceType   *valueObject.ServiceType   `json:"serviceType,omitempty"`
}

type ReadInstalledServicesItemsResponse struct {
	Pagination                   tkDto.Pagination              `json:"pagination"`
	InstalledServices            []entity.InstalledService     `json:"installedServices"`
	InstalledServicesWithMetrics []InstalledServiceWithMetrics `json:"installedServicesWithMetrics"`
}
