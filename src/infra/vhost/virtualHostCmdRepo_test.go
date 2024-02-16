package vhostInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestVirtualHostCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	t.Run("AddAlias", func(t *testing.T) {
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

	t.Run("AddTopLevel", func(t *testing.T) {
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

	t.Run("CreateMapping", func(t *testing.T) {
		responseCode := valueObject.NewHttpResponseCodePanic(403)

		createDto := dto.NewCreateMapping(
			valueObject.NewFqdnPanic("speedia.org"),
			valueObject.NewMappingPathPanic("/"),
			valueObject.NewMappingMatchPatternPanic("begins-with"),
			valueObject.NewMappingTargetTypePanic("response-code"),
			nil,
			nil,
			&responseCode,
		)

		err := VirtualHostCmdRepo{}.CreateMapping(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteMapping", func(t *testing.T) {
		hostname := valueObject.NewFqdnPanic("speedia.org")
		responseCode := valueObject.NewHttpResponseCodePanic(403)
		mapping := entity.NewMapping(
			valueObject.NewMappingIdPanic(0),
			hostname,
			valueObject.NewMappingPathPanic("/"),
			valueObject.NewMappingMatchPatternPanic("begins-with"),
			valueObject.NewMappingTargetTypePanic("response-code"),
			nil,
			nil,
			&responseCode,
		)

		err := VirtualHostCmdRepo{}.DeleteMapping(mapping)
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
