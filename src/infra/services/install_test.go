package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func TestInstall(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("InstallOls", func(t *testing.T) {
		t.Skip("Skip ols install test")
		err := Install(
			valueObject.NewServiceNamePanic("openlitespeed"),
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

	t.Run("InstallRedis", func(t *testing.T) {
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
