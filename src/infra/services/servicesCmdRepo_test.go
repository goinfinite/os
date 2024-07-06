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

		portBinding, err := valueObject.NewPortBindingFromString(
			"8000/http",
		)
		if err != nil {
			t.Errorf("NewPortBindingFromStringFailed : %v", err)
			return
		}

		createDto := dto.NewCreateCustomService(
			valueObject.NewServiceNamePanic("python-ws"),
			valueObject.NewServiceTypePanic("webserver"),
			valueObject.NewUnixCommandPanic("python3 -m http.server"),
			[]valueObject.ServiceEnv{},
			[]valueObject.PortBinding{portBinding},
			true,
			nil,
		)

		err = servicesCmdRepo.CreateCustom(createDto)
		if err != nil {
			t.Errorf("CreateCustomServiceFailed : %v", err)
			return
		}
	})
}
