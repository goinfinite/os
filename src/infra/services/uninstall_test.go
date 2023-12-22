package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestUninstall(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("UninstallPhp", func(t *testing.T) {
		t.Skip("SkipPhpUninstallTest")
		err := Uninstall(
			valueObject.NewServiceNamePanic("php"),
		)
		if err != nil {
			t.Errorf("PhpUninstallFailed : %v", err)
			return
		}
	})

	t.Run("UninstallMariadb", func(t *testing.T) {
		t.Skip("SkipMariadbUninstallTest")
		err := Uninstall(
			valueObject.NewServiceNamePanic("mariadb"),
		)
		if err != nil {
			t.Errorf("MariadbUninstallFailed : %v", err)
			return
		}
	})

	t.Run("UninstallNode", func(t *testing.T) {
		t.Skip("SkipNodeUninstallTest")
		err := Uninstall(
			valueObject.NewServiceNamePanic("node"),
		)
		if err != nil {
			t.Errorf("NodeUninstallFailed : %v", err)
			return
		}
	})

	t.Run("UninstallRedis", func(t *testing.T) {
		t.Skip("SkipRedisUninstallTest")
		err := Uninstall(
			valueObject.NewServiceNamePanic("redis"),
		)
		if err != nil {
			t.Errorf("RedisUninstallFailed : %v", err)
			return
		}
	})
}
