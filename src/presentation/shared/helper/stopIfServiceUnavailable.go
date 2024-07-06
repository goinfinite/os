package sharedHelper

import (
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func StopIfServiceUnavailable(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	serviceName valueObject.ServiceName,
) {
	isServiceRunning := true

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(persistentDbSvc)
	availableSvc, err := servicesQueryRepo.GetByName(serviceName)
	if err != nil {
		isServiceRunning = false
	}

	if availableSvc.Status.String() != "running" {
		isServiceRunning = false
	}

	if !isServiceRunning {
		panic("ServiceUnavailable: " + serviceName.String())
	}
}
