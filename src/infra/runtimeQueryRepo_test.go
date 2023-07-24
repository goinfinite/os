package infra

import (
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/valueObject"
	servicesInfra "github.com/speedianet/sam/src/infra/services"
)

func TestRuntimeQueryRepo(t *testing.T) {
	t.Skip("SkipRuntimeQueryRepoTest")
	testHelpers.LoadEnvVars()

	servicesInfra.Install(
		valueObject.NewServiceNamePanic("openlitespeed"),
		nil,
	)

	t.Run("ReturnPhpVersionsList", func(t *testing.T) {
		repo := RuntimeQueryRepo{}
		phpVersions, err := repo.GetPhpVersions()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(phpVersions) == 0 {
			t.Errorf("Expected a list of php versions, got %v", phpVersions)
		}
	})
}
