package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestInstall(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("InstallPhp", func(t *testing.T) {
		t.Skip("Skip php install test")
		version, _ := valueObject.NewServiceVersion("7.4")
		err := Install(
			valueObject.NewServiceNamePanic("php"),
			&version,
		)
		if err != nil {
			t.Errorf("Install() error = %v", err)
			return
		}
	})

	t.Run("InstallNode", func(t *testing.T) {
		t.Skip("Skip node install test")
		err := Install(
			valueObject.NewServiceNamePanic("node"),
			nil,
		)
		if err != nil {
			t.Errorf("Install() error = %v", err)
			return
		}
	})

	t.Run("InstallMysql", func(t *testing.T) {
		t.Skip("Skip mysql install test")
		err := Install(
			valueObject.NewServiceNamePanic("mysql"),
			nil,
		)
		if err != nil {
			t.Errorf("Install() error = %v", err)
			return
		}
	})

	t.Run("InstallRedis", func(t *testing.T) {
		t.Skip("Skip redis install test")
		version, _ := valueObject.NewServiceVersion("7.0")
		err := Install(
			valueObject.NewServiceNamePanic("redis"),
			&version,
		)
		if err != nil {
			t.Errorf("Install() error = %v", err)
			return
		}
	})
}
