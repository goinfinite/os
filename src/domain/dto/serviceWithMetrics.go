package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServiceWithMetrics struct {
	entity.Service
	Metrics *valueObject.ServiceMetrics `json:"metrics,omitempty"`
}

func NewServiceWithMetrics(
	service entity.Service,
	metrics *valueObject.ServiceMetrics,
) ServiceWithMetrics {
	return ServiceWithMetrics{
		Service: service,
		Metrics: metrics,
	}
}
