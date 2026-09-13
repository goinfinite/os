package liaisonHelper

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

const ServiceUnavailableError = "ServiceUnavailable"

type ServiceAvailabilityInspector struct {
	servicesQueryRepo *servicesInfra.ServicesQueryRepo
}

func NewServiceAvailabilityInspector(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServiceAvailabilityInspector {
	return &ServiceAvailabilityInspector{
		servicesQueryRepo: servicesInfra.NewServicesQueryRepo(persistentDbSvc),
	}
}

func (inspector *ServiceAvailabilityInspector) IsAvailable(
	serviceName valueObject.ServiceName,
) bool {
	availableService, err := inspector.servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{ServiceName: &serviceName},
	)
	if err != nil {
		return false
	}

	return availableService.Status == valueObject.ServiceStatusRunning
}

func NewServiceUnavailableResponse() tkPresentation.LiaisonResponse {
	return tkPresentation.NewLiaisonResponseNoMessage(
		tkPresentation.LiaisonResponseStatusServiceUnavailable,
		ServiceUnavailableError,
	)
}
