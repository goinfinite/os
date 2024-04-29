package vhostInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestVirtualHostCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	t.Run("CreateAlias", func(t *testing.T) {
		parentDomain := valueObject.NewFqdnPanic("speedia.net")

		createDto := dto.NewCreateVirtualHost(
			valueObject.NewFqdnPanic("speedia.com"),
			valueObject.NewVirtualHostTypePanic("alias"),
			&parentDomain,
		)

		err := VirtualHostCmdRepo{}.Create(createDto)

		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("CreateTopLevel", func(t *testing.T) {
		createDto := dto.NewCreateVirtualHost(
			valueObject.NewFqdnPanic("speedia.org"),
			valueObject.NewVirtualHostTypePanic("top-level"),
			nil,
		)

		err := VirtualHostCmdRepo{}.Create(createDto)

		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteTopLevelAndAliases", func(t *testing.T) {
		hostnames := []valueObject.Fqdn{
			valueObject.NewFqdnPanic("speedia.com"),
			valueObject.NewFqdnPanic("speedia.org"),
		}

		for _, hostname := range hostnames {
			vhostEntity, err := VirtualHostQueryRepo{}.GetByHostname(hostname)
			if err != nil {
				t.Errorf("ExpectedNoErrorButGot: %v", err)
			}

			err = VirtualHostCmdRepo{}.Delete(vhostEntity)
			if err != nil {
				t.Errorf("ExpectedNoErrorButGot: %v", err)
			}
		}
	})
}
