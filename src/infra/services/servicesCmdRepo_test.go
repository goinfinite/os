package servicesInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

func TestServiceCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	servicesCmdRepo := NewServicesCmdRepo(persistentDbSvc)

	t.Run("CreateCustomService", func(t *testing.T) {
		t.Skip("SkipCreateCustomServiceTest")

		portBinding, err := valueObject.NewPortBinding(
			"8000/http",
		)
		if err != nil {
			t.Errorf("NewPortBindingFailed : %v", err)
			return
		}

		serviceName, _ := valueObject.NewServiceName("python-ws")
		serviceType, _ := valueObject.NewServiceType("webserver")
		unixCommand, _ := valueObject.NewUnixCommand("python3 -m http.server")
		err = servicesCmdRepo.CreateCustom(dto.CreateCustomService{
			Name:              serviceName,
			Type:              serviceType,
			StartCmd:          unixCommand,
			PortBindings:      []valueObject.PortBinding{portBinding},
			OperatorAccountId: valueObject.AccountIdSystem,
			OperatorIpAddress: valueObject.IpAddressSystem,
		})
		if err != nil {
			t.Errorf("CreateCustomServiceFailed : %v", err)
			return
		}
	})
}
