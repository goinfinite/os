package sharedHelper

import (
	"slices"

	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func CheckServiceAvailability(svcName string, requiredSvcNames *[]string) {
	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	availableSvcs, err := servicesQueryRepo.Get()
	if err != nil {
		panic("FailedToGetAvailableServices: " + err.Error())
	}

	serviceIsRunning := false
	for _, availableSvc := range availableSvcs {
		availableSvcName := availableSvc.Name.String()

		hasRequiredSvcNames := requiredSvcNames != nil
		if hasRequiredSvcNames && !slices.Contains(*requiredSvcNames, availableSvcName) {
			continue
		}

		if availableSvcName != svcName {
			continue
		}

		if availableSvc.Status.String() != "running" {
			continue
		}

		serviceIsRunning = true
		break
	}

	if !serviceIsRunning {
		panic("ServiceUnavailable: " + svcName)
	}
}
