package runtimeInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

func TestRuntimeQueryRepo(t *testing.T) {
	t.Skip("SkipRuntimeQueryRepoTest")
	testHelpers.LoadEnvVars()

	repo := RuntimeQueryRepo{}

	t.Run("ReturnPhpVersionsList", func(t *testing.T) {
		phpVersions, err := repo.ReadPhpVersionsInstalled()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(phpVersions) == 0 {
			t.Errorf("Expected a list of php versions, got %v", phpVersions)
		}
	})

	t.Run("ReturnPhpConfigs", func(t *testing.T) {
		primaryVhost, err := infraHelper.ReadPrimaryVirtualHostHostname()
		if err != nil {
			t.Errorf("PrimaryVirtualHostNotFound")
		}

		hostname, _ := valueObject.NewFqdn(primaryVhost.String())
		phpConfigs, err := repo.ReadPhpConfigs(hostname)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(phpConfigs.Modules) == 0 {
			t.Errorf("Expected a list of php modules, got %v", phpConfigs)
		}
	})
}
