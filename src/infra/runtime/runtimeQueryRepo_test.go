package runtimeInfra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/valueObject"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func TestRuntimeQueryRepo(t *testing.T) {
	t.Skip("SkipRuntimeQueryRepoTest")
	testHelpers.LoadEnvVars()

	err := servicesInfra.AddInstallableSimplified("php")
	if err != nil {
		t.Errorf("InstallDependenciesFail: %v", err)
		return
	}

	repo := RuntimeQueryRepo{}

	t.Run("ReturnPhpVersionsList", func(t *testing.T) {
		phpVersions, err := repo.GetPhpVersionsInstalled()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(phpVersions) == 0 {
			t.Errorf("Expected a list of php versions, got %v", phpVersions)
		}
	})

	t.Run("ReturnPhpConfigs", func(t *testing.T) {
		hostname := valueObject.NewFqdnPanic(os.Getenv("VIRTUAL_HOST"))
		phpConfigs, err := repo.GetPhpConfigs(hostname)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(phpConfigs.Modules) == 0 {
			t.Errorf("Expected a list of php modules, got %v", phpConfigs)
		}
	})
}
