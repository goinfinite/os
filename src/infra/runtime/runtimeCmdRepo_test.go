package runtimeInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func TestRuntimeCmdRepo(t *testing.T) {
	t.Skip("SkipRuntimeCmdRepoTest")
	testHelpers.LoadEnvVars()

	err := servicesInfra.AddInstallableSimplified("php")
	if err != nil {
		t.Errorf("InstallDependenciesFail: %v", err)
		return
	}

	t.Run("UpdatePhpVersion", func(t *testing.T) {
		primaryHostname, err := infraHelper.GetPrimaryHostname()
		if err != nil {
			t.Errorf("PrimaryHostnameNotFound")
		}

		err = RuntimeCmdRepo{}.UpdatePhpVersion(
			valueObject.NewFqdnPanic(primaryHostname.String()),
			valueObject.NewPhpVersionPanic("8.1"),
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("UpdatePhpSettings", func(t *testing.T) {
		primaryHostname, err := infraHelper.GetPrimaryHostname()
		if err != nil {
			t.Errorf("PrimaryHostnameNotFound")
		}

		err = RuntimeCmdRepo{}.UpdatePhpSettings(
			valueObject.NewFqdnPanic(primaryHostname.String()),
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
		primaryHostname, err := infraHelper.GetPrimaryHostname()
		if err != nil {
			t.Errorf("PrimaryHostnameNotFound")
		}

		err = RuntimeCmdRepo{}.UpdatePhpModules(
			valueObject.NewFqdnPanic(primaryHostname.String()),
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

		err = RuntimeCmdRepo{}.UpdatePhpModules(
			valueObject.NewFqdnPanic(primaryHostname.String()),
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
