package infra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestVirtualHostCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	t.Run("AddAlias", func(t *testing.T) {
		parentDomain := valueObject.NewFqdnPanic("speedia.net")

		addDto := dto.NewAddVirtualHost(
			valueObject.NewFqdnPanic("speedia.com"),
			valueObject.NewVirtualHostTypePanic("alias"),
			&parentDomain,
		)

		err := VirtualHostCmdRepo{}.Add(addDto)

		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("AddTopLevel", func(t *testing.T) {
		addDto := dto.NewAddVirtualHost(
			valueObject.NewFqdnPanic("speedia.org"),
			valueObject.NewVirtualHostTypePanic("top-level"),
			nil,
		)

		err := VirtualHostCmdRepo{}.Add(addDto)

		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
