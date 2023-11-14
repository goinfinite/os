package infra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/valueObject"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func TestWsQueryRepo(t *testing.T) {
	t.Skip("SkipWsQueryRepoTest")
	testHelpers.LoadEnvVars()

	servicesInfra.Install(
		valueObject.NewServiceNamePanic("openlitespeed"),
		nil,
	)

	t.Run("ReturnVirtualHostsList", func(t *testing.T) {
		repo := WsQueryRepo{}
		vhosts, err := repo.GetVirtualHosts()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(vhosts) == 0 {
			t.Errorf("Expected a list of vhosts, got %v", vhosts)
		}
	})
}
