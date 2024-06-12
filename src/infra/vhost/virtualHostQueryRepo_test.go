package vhostInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

func TestVirtualHostQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	vhostQueryRepo := NewVirtualHostQueryRepo(persistentDbSvc)

	t.Run("Read", func(t *testing.T) {
		vhosts, err := vhostQueryRepo.Read()
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		if len(vhosts) > 0 {
			t.Errorf("ExpectingEmptySliceButGot: %v", len(vhosts))
		}
	})

	t.Run("ReadByHostname", func(t *testing.T) {
		hostname, _ := infraHelper.GetPrimaryVirtualHost()
		vhost, err := vhostQueryRepo.ReadByHostname(hostname)
		if err != nil && err.Error() != "VhostNotFound" {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		if vhost.Hostname != "" {
			t.Errorf("ExpectingEmptyVhostButGot: %v", vhost)
		}
	})

	t.Run("ReadAliasesByParentHostname", func(t *testing.T) {
		hostname, _ := infraHelper.GetPrimaryVirtualHost()
		aliases, err := vhostQueryRepo.ReadAliasesByParentHostname(hostname)
		if err != nil && err.Error() != "VhostNotFound" {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		if len(aliases) > 0 {
			t.Errorf("ExpectingEmptySliceButGot: %v", len(aliases))
		}
	})
}
