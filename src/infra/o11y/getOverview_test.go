package o11yInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestGetOverview(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetOverview", func(t *testing.T) {
		getOverviewRepo := GetOverview{}
		_, err := getOverviewRepo.Get()
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
