package o11yInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	internalDatabaseInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

func TestGetOverview(t *testing.T) {
	testHelpers.LoadEnvVars()

	transientDbSvc, err := internalDatabaseInfra.NewTransientDatabaseService()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	o11yQueryRepo := NewO11yQueryRepo(transientDbSvc)

	t.Run("GetOverview", func(t *testing.T) {
		_, err := o11yQueryRepo.GetOverview()
		if err != nil {
			t.Errorf("Expected nil, got %s", err.Error())
		}
	})
}
