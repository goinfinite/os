package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func TestUninstall(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("UninstallOls", func(t *testing.T) {
		t.Skip("Skip ols uninstall test")
		err := Uninstall(
			valueObject.NewServiceNamePanic("openlitespeed"),
		)
		if err != nil {
			t.Errorf("Uninstall() error = %v", err)
			return
		}
	})

	t.Run("UninstallMysql", func(t *testing.T) {
		t.Skip("Skip mysql uninstall test")
		err := Uninstall(
			valueObject.NewServiceNamePanic("mysql"),
		)
		if err != nil {
			t.Errorf("Uninstall() error = %v", err)
			return
		}
	})

	t.Run("UninstallNode", func(t *testing.T) {
		t.Skip("Skip node uninstall test")
		err := Uninstall(
			valueObject.NewServiceNamePanic("node"),
		)
		if err != nil {
			t.Errorf("Uninstall() error = %v", err)
			return
		}
	})

	t.Run("UninstallRedis", func(t *testing.T) {
		t.Skip("Skip redis uninstall test")
		err := Uninstall(
			valueObject.NewServiceNamePanic("redis"),
		)
		if err != nil {
			t.Errorf("Uninstall() error = %v", err)
			return
		}
	})
}
