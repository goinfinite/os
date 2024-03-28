package vhostInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

func TestVirtualHostQueryRepo(t *testing.T) {
	vhostQueryRepo := VirtualHostQueryRepo{}
	testHelpers.LoadEnvVars()

	t.Run("GetVirtualHosts", func(t *testing.T) {
		vhosts, err := vhostQueryRepo.Get()
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		if len(vhosts) == 0 {
			t.Errorf("ExpectingNonEmptySliceButGot: %v", vhosts)
		}
	})

	t.Run("GetVirtualHostsWithMappings", func(t *testing.T) {
		dummyMapping := `
location / {
	return 301 https://speedia.net;
}
`

		err := infraHelper.UpdateFile(
			"/app/conf/nginx/mapping/primary.conf",
			dummyMapping,
			true,
		)
		if err != nil {
			t.Errorf("UpdateMappingFileError: %s", err.Error())
		}

		vhosts, err := vhostQueryRepo.GetWithMappings()
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		if len(vhosts) == 0 {
			t.Errorf("ExpectingNonEmptySliceButGot: %v", vhosts)
		}
	})

	t.Run("TestDomainNotMappedToServer", func(t *testing.T) {
		primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
		if err != nil {
			t.Errorf("GetPrimaryVirtualHostError: %s", err.Error())
		}
		randomHash := valueObject.NewHashPanic("randomhash")

		isDomainMapped := vhostQueryRepo.IsDomainMappedToServer(
			primaryVhost,
			randomHash,
		)

		if isDomainMapped {
			t.Errorf("ExpectingFalseButGot: %v", isDomainMapped)
		}
	})
}
