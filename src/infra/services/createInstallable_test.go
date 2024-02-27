package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestCreateInstallable(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("InstallPhp", func(t *testing.T) {
		t.Skip("SkipPhpInstallTest")

		err := CreateInstallableSimplified("php")
		if err != nil {
			t.Errorf("PhpInstallFailed : %v", err)
			return
		}
	})

	t.Run("InstallNode", func(t *testing.T) {
		t.Skip("SkipNodeInstallTest")

		err := CreateInstallableSimplified("node")
		if err != nil {
			t.Errorf("NodeInstallFailed : %v", err)
			return
		}
	})

	t.Run("InstallMariaDb", func(t *testing.T) {
		t.Skip("SkipMariaDbInstallTest")

		err := CreateInstallableSimplified("mariadb")
		if err != nil {
			t.Errorf("MariaDbInstallFailed : %v", err)
			return
		}
	})

	t.Run("InstallRedis", func(t *testing.T) {
		t.Skip("SkipRedisInstallTest")

		err := CreateInstallableSimplified("redis")
		if err != nil {
			t.Errorf("RedisInstallFailed : %v", err)
			return
		}
	})
}
