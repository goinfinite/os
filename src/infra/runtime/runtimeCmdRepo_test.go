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

	t.Run("UpdatePhpVersion", func(t *testing.T) {
		primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
		if err != nil {
			t.Errorf("PrimaryVirtualHostNotFound")
		}

		err = runtimeCmdRepo.UpdatePhpVersion(
			valueObject.NewFqdnPanic(primaryVhost.String()),
			valueObject.NewPhpVersionPanic("8.1"),
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("UpdatePhpSettings", func(t *testing.T) {
		primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
		if err != nil {
			t.Errorf("PrimaryVirtualHostNotFound")
		}

		err = runtimeCmdRepo.UpdatePhpSettings(
			valueObject.NewFqdnPanic(primaryVhost.String()),
			[]entity.PhpSetting{
				entity.NewPhpSetting(
					valueObject.NewPhpSettingNamePanic("display_errors"),
					valueObject.NewPhpSettingValuePanic("Off"),
					nil,
				),
			},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("UpdatePhpModules", func(t *testing.T) {
		primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
		if err != nil {
			t.Errorf("PrimaryVirtualHostNotFound")
		}

		err = runtimeCmdRepo.UpdatePhpModules(
			valueObject.NewFqdnPanic(primaryVhost.String()),
			[]entity.PhpModule{
				entity.NewPhpModule(
					valueObject.NewPhpModuleNamePanic("ioncube"),
					true,
				),
			},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpModules(
			valueObject.NewFqdnPanic(primaryVhost.String()),
			[]entity.PhpModule{
				entity.NewPhpModule(
					valueObject.NewPhpModuleNamePanic("ioncube"),
					false,
				),
			},
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
