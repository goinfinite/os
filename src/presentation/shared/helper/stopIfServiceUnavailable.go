package sharedHelper

import (
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func StopIfServiceUnavailable(svcName string) {
	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	availableSvcs, err := servicesQueryRepo.Get()
	if err != nil {
		panic("FailedToGetAvailableServices: " + err.Error())
	}

	for _, availableSvc := range availableSvcs {
		availableSvcName := availableSvc.Name.String()
		if availableSvcName != svcName {
			continue
		}

		if availableSvc.Status.String() != "running" {
			panic("ServiceUnavailable: " + svcName)
		}
	}
}
