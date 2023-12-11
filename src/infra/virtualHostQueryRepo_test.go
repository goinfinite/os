package infra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestVirtualHostQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetVirtualHosts", func(t *testing.T) {
		repo := VirtualHostQueryRepo{}
		vhosts, err := repo.Get()

		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		if len(vhosts) == 0 {
			t.Errorf("ExpectingNonEmptySliceButGot: %v", vhosts)
		}
	})
}
