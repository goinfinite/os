package vhostInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
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
}
