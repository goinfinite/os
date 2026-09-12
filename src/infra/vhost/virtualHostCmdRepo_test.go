package vhostInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestVirtualHostCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	vhostCmdRepo := NewVirtualHostCmdRepo(persistentDbSvc)
	vhostQueryRepo := NewVirtualHostQueryRepo(persistentDbSvc)

	vhostName, _ := NewVirtualHostHelpers().ReadPrimaryVirtualHostHostname()

	// Note: Setup/teardown are intentionally inline — test independence
	// requires each file to own its preconditions, even if it duplicates code.
	existingVhosts, err := vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
		Pagination: tkDto.PaginationSingleItem,
		Hostname:   &vhostName,
	})
	if err != nil {
		t.Fatalf("VirtualHostPreconditionReadFailed: %v", err)
	}
	if len(existingVhosts.VirtualHosts) > 0 {
		err = vhostCmdRepo.Delete(vhostName)
		if err != nil {
			t.Fatalf("VirtualHostPreconditionDeleteFailed: %v", err)
		}
	}

	t.Run("Create", func(t *testing.T) {
		vhostType, _ := valueObject.NewVirtualHostType("top-level")
		operatorAccountId, _ := tkValueObject.NewAccountId(0)
		ipAddress := tkValueObject.IpAddressLocal

		err := vhostCmdRepo.Create(dto.NewCreateVirtualHost(
			vhostName, vhostType, nil, nil, operatorAccountId, ipAddress,
		))
		if err != nil {
			t.Fatalf("ExpectingNoErrorButGot: %v", err)
		}

		readResponse, readErr := vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
			Pagination: tkDto.PaginationSingleItem,
			Hostname:   &vhostName,
		})
		if readErr != nil {
			t.Fatalf("VirtualHostReadFailed: %v", readErr)
		}
		if len(readResponse.VirtualHosts) == 0 {
			t.Errorf("CreatedVirtualHostNotFound: %s", vhostName.String())
		}
	})

	t.Run("Update", func(t *testing.T) {
		isWildcard := true
		err := vhostCmdRepo.Update(dto.UpdateVirtualHost{
			Hostname:   vhostName,
			IsWildcard: &isWildcard,
		})
		if err != nil {
			t.Fatalf("ExpectingNoErrorButGot: %v", err)
		}

		readResponse, readErr := vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
			Pagination: tkDto.PaginationSingleItem,
			Hostname:   &vhostName,
		})
		if readErr != nil {
			t.Fatalf("VirtualHostReadFailed: %v", readErr)
		}
		if len(readResponse.VirtualHosts) == 0 {
			t.Fatalf("UpdatedVirtualHostNotFound: %s", vhostName.String())
		}
		if !readResponse.VirtualHosts[0].IsWildcard {
			t.Errorf("VirtualHostShouldBeWildcard")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := vhostCmdRepo.Delete(vhostName)
		if err != nil {
			t.Fatalf("ExpectingNoErrorButGot: %v", err)
		}

		readResponse, readErr := vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
			Pagination: tkDto.PaginationSingleItem,
			Hostname:   &vhostName,
		})
		if readErr != nil {
			t.Fatalf("VirtualHostReadFailed: %v", readErr)
		}
		if len(readResponse.VirtualHosts) > 0 {
			t.Errorf("DeletedVirtualHostStillFound: %s", vhostName.String())
		}
	})
}
