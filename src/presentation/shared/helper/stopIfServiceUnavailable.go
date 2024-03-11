package sharedHelper

import (
	"github.com/speedianet/os/src/domain/valueObject"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func StopIfServiceUnavailable(svcNameStr string) {
	svcName := valueObject.NewServiceNamePanic(svcNameStr)

	isServiceRunning := true

	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	availableSvc, err := servicesQueryRepo.GetByName(svcName)
	if err != nil {
		isServiceRunning = false
	}

	if availableSvc.Status.String() != "running" {
		isServiceRunning = false
	}

	if !isServiceRunning {
		panic("ServiceUnavailable: " + svcName.String())
	}
}
