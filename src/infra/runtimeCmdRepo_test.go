package infra

import (
	"os"
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/valueObject"
	servicesInfra "github.com/speedianet/sam/src/infra/services"
)

func TestRuntimeCmdRepo(t *testing.T) {
	t.Skip("SkipRuntimeCmdRepoTest")
	testHelpers.LoadEnvVars()

	servicesInfra.Install(
		valueObject.NewServiceNamePanic("openlitespeed"),
		nil,
	)

	t.Run("UpdatePhpVersion", func(t *testing.T) {
		err := RuntimeCmdRepo{}.UpdatePhpVersion(
			valueObject.NewFqdnPanic(os.Getenv("VIRTUAL_HOST")),
			valueObject.NewPhpVersionPanic("8.1"),
		)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
