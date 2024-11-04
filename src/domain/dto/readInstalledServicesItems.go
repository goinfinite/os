package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ReadInstalledServicesItemsRequest struct {
	Pagination           Pagination                 `json:"pagination"`
	Name                 *valueObject.ServiceName   `json:"name,omitempty"`
	Nature               *valueObject.ServiceNature `json:"nature,omitempty"`
	Type                 *valueObject.ServiceType   `json:"type,omitempty"`
	ShouldIncludeMetrics bool                       `json:"shouldIncludeMetrics"`
}

type ReadInstalledServicesItemsResponse struct {
	Pagination Pagination                    `json:"pagination"`
	Items      []InstalledServiceWithMetrics `json:"items"`
}
