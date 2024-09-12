package sharedHelper

import (
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

const ServiceUnavailableError = "ServiceUnavailable"

func IsServiceAvailable(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	serviceName valueObject.ServiceName,
) bool {
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(persistentDbSvc)
	availableSvc, err := servicesQueryRepo.ReadByName(serviceName)
	if err != nil {
		return false
	}

	return availableSvc.Status.String() == "running"
}
