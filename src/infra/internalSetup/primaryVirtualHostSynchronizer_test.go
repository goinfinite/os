package internalSetupInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
)

func TestPrimaryVirtualHostSynchronizer(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc := testHelpers.GetPersistentDbSvc()
	sync := NewPrimaryVirtualHostSynchronizer(persistentDbSvc)

	t.Run("PhpConfUpdaterSkipsWhenPhpUninstalled", func(t *testing.T) {
		err := sync.phpConfUpdater()
		if err != nil {
			t.Errorf("PhpConfUpdaterShouldReturnNilWhenPhpUninstalled: %v", err)
		}
	})
}
