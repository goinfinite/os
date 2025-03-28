package runtimeInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

func TestRuntimeCmdRepo(t *testing.T) {
	t.Skip("SkipRuntimeCmdRepoTest")
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	runtimeCmdRepo := NewRuntimeCmdRepo(persistentDbSvc)

	primaryVhost, _ := infraHelper.ReadPrimaryVirtualHostHostname()
	phpVersion, _ := valueObject.NewPhpVersion("8.1")

	t.Run("UpdatePhpVersion", func(t *testing.T) {
		err := runtimeCmdRepo.UpdatePhpVersion(primaryVhost, phpVersion)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("UpdatePhpSettings", func(t *testing.T) {
		phpSettingName, _ := valueObject.NewPhpSettingName("display_errors")
		phpSettingType, _ := valueObject.NewPhpSettingType("select")
		phpSettingValue, _ := valueObject.NewPhpSettingValue("Off")

		err := runtimeCmdRepo.UpdatePhpSettings(
			primaryVhost,
			[]entity.PhpSetting{
				entity.NewPhpSetting(
					phpSettingName, phpSettingType, phpSettingValue, nil,
				),
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
