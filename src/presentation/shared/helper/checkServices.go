package sharedHelper

import (
	"errors"

	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
)

func CheckServices(serviceNameStr string) error {
	servicesQueryRepo := infra.ServicesQueryRepo{}

	serviceName, err := valueObject.NewServiceName(serviceNameStr)
	if err != nil {
		return err
	}

	currentSvcStatus, err := servicesQueryRepo.GetByName(serviceName)
	if err != nil {
		return err
	}

	isRunning := currentSvcStatus.Status.String() == "running"
	if !isRunning {
		return errors.New("ServiceUnavailableError")
	}

	return nil
}
