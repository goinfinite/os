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

	primaryVhost, _ := infraHelper.GetPrimaryVirtualHost()
	phpVersion, _ := valueObject.NewPhpVersion("8.1")

	t.Run("UpdatePhpVersion", func(t *testing.T) {
		err := runtimeCmdRepo.UpdatePhpVersion(primaryVhost, phpVersion)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("UpdatePhpSettings", func(t *testing.T) {
		phpSettingName, _ := valueObject.NewPhpSettingName("display_errors")
		phpSettingValue, _ := valueObject.NewPhpSettingValue("Off")

		err := runtimeCmdRepo.UpdatePhpSettings(
			primaryVhost,
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

		err := runtimeCmdRepo.UpdatePhpModules(
			primaryVhost,
			[]entity.PhpModule{entity.NewPhpModule(phpModuleName, true)},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpModules(
			primaryVhost,
			[]entity.PhpModule{entity.NewPhpModule(phpModuleName, false)},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
