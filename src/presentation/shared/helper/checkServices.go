package sharedHelper

import (
	"errors"
	"slices"

	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func CheckServices(servicesNames []string) error {
	servicesQueryRepo := servicesInfra.ServicesQueryRepo{}
	services, err := servicesQueryRepo.Get()
	if err != nil {
		return err
	}

	serviceIsRunning := false
	for _, service := range services {
		if service.Status.String() != "running" {
			continue
		}

		if !slices.Contains(servicesNames, service.Name.String()) {
			continue
		}

		serviceIsRunning = true
	}

	if !serviceIsRunning {
		return errors.New("ServiceUnavailable")
	}

	return nil
}
