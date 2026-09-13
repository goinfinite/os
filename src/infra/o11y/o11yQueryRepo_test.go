package o11yInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
)

func TestGetOverview(t *testing.T) {
	testHelpers.LoadEnvVars()

	transientDbSvc, err := tkInfraDb.NewTransientDatabaseService()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	o11yQueryRepo := NewO11yQueryRepo(transientDbSvc)

	t.Run("GetOverview", func(t *testing.T) {
		_, err := o11yQueryRepo.ReadOverview(true)
		if err != nil {
			t.Errorf("Expected nil, got %s", err.Error())
		}
	})
}
