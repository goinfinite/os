package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
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

		createDto := dto.NewCreateCustomService(
			serviceName, serviceType, unixCommand, []valueObject.ServiceEnv{},
			[]valueObject.PortBinding{portBinding}, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)

		err = servicesCmdRepo.CreateCustom(createDto)
		if err != nil {
			t.Errorf("CreateCustomServiceFailed : %v", err)
			return
		}
	})
}
