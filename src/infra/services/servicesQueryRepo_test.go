package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

func TestServicesQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	servicesQueryRepo := NewServicesQueryRepo(persistentDbSvc)

	t.Run("ReturnServicesList", func(t *testing.T) {
		services, err := servicesQueryRepo.Read()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(services) == 0 {
			t.Errorf("Expected a list of services, got %v", services)
		}
	})
}
