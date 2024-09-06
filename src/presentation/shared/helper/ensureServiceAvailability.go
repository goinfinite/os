package sharedHelper

import (
	"errors"

	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func EnsureServiceAvailability(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	serviceName valueObject.ServiceName,
) error {
	isServiceRunning := true

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(persistentDbSvc)
	availableSvc, err := servicesQueryRepo.ReadByName(serviceName)
	if err != nil {
		isServiceRunning = false
	}

	if availableSvc.Status.String() != "running" {
		isServiceRunning = false
	}

	if !isServiceRunning {
		return errors.New("ServiceUnavailable: " + serviceName.String())
	}

	return nil
}
