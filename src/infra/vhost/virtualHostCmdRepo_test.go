package vhostInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
)

func TestVirtualHostCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	vhostCmdRepo := NewVirtualHostCmdRepo(persistentDbSvc)
	vhostQueryRepo := NewVirtualHostQueryRepo(persistentDbSvc)

	vhostName, _ := infraHelper.GetPrimaryVirtualHost()

	t.Run("Create", func(t *testing.T) {
		vhostType := valueObject.NewVirtualHostTypePanic("top-level")
		createDto := dto.NewCreateVirtualHost(vhostName, vhostType, nil)

		err := vhostCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		vhost, err := vhostQueryRepo.ReadByHostname(vhostName)
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}

		err = vhostCmdRepo.Delete(vhost)
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}
	})
}
