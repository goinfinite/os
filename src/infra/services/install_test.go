package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func TestInstall(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("InstallOLS", func(t *testing.T) {
		err := Install(
			valueObject.NewServiceNamePanic("openlitespeed"),
			nil,
		)
		if err != nil {
			t.Errorf("Install() error = %v", err)
			return
		}
	})
}
