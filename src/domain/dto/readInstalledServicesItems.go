package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadInstalledServicesItemsRequest struct {
	Pagination           Pagination                 `json:"pagination"`
	ServiceName          *valueObject.ServiceName   `json:"serviceName,omitempty"`
	ServiceNature        *valueObject.ServiceNature `json:"serviceNature,omitempty"`
	ServiceStatus        *valueObject.ServiceStatus `json:"serviceStatus,omitempty"`
	ServiceType          *valueObject.ServiceType   `json:"serviceType,omitempty"`
	ShouldIncludeMetrics *bool                      `json:"shouldIncludeMetrics,omitempty"`
}

type ReadFirstInstalledServiceItemsRequest struct {
	Pagination    Pagination                 `json:"pagination"`
	ServiceName   *valueObject.ServiceName   `json:"serviceName,omitempty"`
	ServiceNature *valueObject.ServiceNature `json:"serviceNature,omitempty"`
	ServiceStatus *valueObject.ServiceStatus `json:"serviceStatus,omitempty"`
	ServiceType   *valueObject.ServiceType   `json:"serviceType,omitempty"`
}

type ReadInstalledServicesItemsResponse struct {
	Pagination                   Pagination                    `json:"pagination"`
	InstalledServices            []entity.InstalledService     `json:"installedServices"`
	InstalledServicesWithMetrics []InstalledServiceWithMetrics `json:"installedServicesWithMetrics"`
}
