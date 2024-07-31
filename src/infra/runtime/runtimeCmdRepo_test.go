package runtimeInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

func TestRuntimeCmdRepo(t *testing.T) {
	t.Skip("SkipRuntimeCmdRepoTest")
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	runtimeCmdRepo := NewRuntimeCmdRepo(persistentDbSvc)

	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		t.Errorf("PrimaryVirtualHostNotFound")
	}
	vhostHostname, _ := valueObject.NewFqdn(primaryVhost.String())

	t.Run("UpdatePhpVersion", func(t *testing.T) {
		phpVersion, _ := valueObject.NewPhpVersion("8.1")

		err := runtimeCmdRepo.UpdatePhpVersion(vhostHostname, phpVersion)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("UpdatePhpSettings", func(t *testing.T) {
		phpSettingName, _ := valueObject.NewPhpSettingName("display_errors")
		phpSettingValue, _ := valueObject.NewPhpSettingValue("Off")

		err = runtimeCmdRepo.UpdatePhpSettings(
			vhostHostname,
			[]entity.PhpSetting{
				entity.NewPhpSetting(phpSettingName, phpSettingValue, nil),
			},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("UpdatePhpModules", func(t *testing.T) {
		phpModuleName, _ := valueObject.NewPhpModuleName("ioncube")

		err = runtimeCmdRepo.UpdatePhpModules(
			vhostHostname,
			[]entity.PhpModule{
				entity.NewPhpModule(phpModuleName, true),
			},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpModules(
			vhostHostname,
			[]entity.PhpModule{
				entity.NewPhpModule(phpModuleName, false),
			},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
