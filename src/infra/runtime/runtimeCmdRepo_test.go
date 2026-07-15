package runtimeInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestRuntimeCmdRepo(t *testing.T) {
	t.Skip("SkipRuntimeCmdRepoTest")
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	runtimeCmdRepo := NewRuntimeCmdRepo(persistentDbSvc)

	primaryVhost, _ := vhostInfra.NewVirtualHostHelpers().
		ReadPrimaryVirtualHostHostname()
	phpVersion, _ := valueObject.NewPhpVersion("8.1")

	t.Run("UpdatePhpVersion", func(t *testing.T) {
		err := runtimeCmdRepo.UpdatePhpVersion(primaryVhost, phpVersion)
		if err != nil {
			t.Errorf("UpdatePhpVersionShouldSucceed: %v", err)
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
			t.Errorf("UpdatePhpSettingsShouldSucceed: %v", err)
		}
	})

	t.Run("UpdatePhpModules", func(t *testing.T) {
		phpModuleName, _ := valueObject.NewPhpModuleName("ioncube")

		err := runtimeCmdRepo.UpdatePhpModules(
			primaryVhost,
			[]entity.PhpModule{entity.NewPhpModule(phpModuleName, true)},
		)
		if err != nil {
			t.Errorf("UpdatePhpModulesEnableShouldSucceed: %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpModules(
			primaryVhost,
			[]entity.PhpModule{entity.NewPhpModule(phpModuleName, false)},
		)
		if err != nil {
			t.Errorf("UpdatePhpModulesDisableShouldSucceed: %v", err)
		}
	})

	t.Run("UpdatePhpVirtualHostHostname", func(t *testing.T) {
		newHostname, _ := tkValueObject.NewFqdn(primaryVhost.String() + ".renamed")

		err := runtimeCmdRepo.UpdatePhpVirtualHostHostname(
			primaryVhost, newHostname, []tkValueObject.Fqdn{},
		)
		if err != nil {
			t.Errorf("UpdatePhpVirtualHostHostnameShouldSucceed: %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpVirtualHostHostname(
			newHostname, primaryVhost, []tkValueObject.Fqdn{},
		)
		if err != nil {
			t.Errorf("UpdatePhpVirtualHostHostnameReverseShouldSucceed: %v", err)
		}
	})

	t.Run("UpdatePhpVirtualHostHostnameNoOp", func(t *testing.T) {
		err := runtimeCmdRepo.UpdatePhpVirtualHostHostname(
			primaryVhost, primaryVhost, []tkValueObject.Fqdn{},
		)
		if err != nil {
			t.Errorf("UpdatePhpVirtualHostHostnameNoOpShouldReturnNil: %v", err)
		}
	})
}
