package internalSetupInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

func TestPrimaryVirtualHostSynchronizer(t *testing.T) {
	t.Skip("SkipPrimaryVirtualHostSynchronizerTest")
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	sync := NewPrimaryVirtualHostSynchronizer(persistentDbSvc)

	t.Run("PhpConfUpdaterSkipsWhenPhpUninstalled", func(t *testing.T) {
		err := sync.phpConfUpdater()
		if err != nil {
			t.Errorf("PhpConfUpdaterShouldReturnNilWhenPhpUninstalled: %v", err)
		}
	})
}
