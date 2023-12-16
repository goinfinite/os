package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestServicesQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	t.Run("ReturnServicesList", func(t *testing.T) {
		repo := ServicesQueryRepo{}
		services, err := repo.Get()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(services) == 0 {
			t.Errorf("Expected a list of services, got %v", services)
		}
	})
}
