package vhostInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

func TestVirtualHostQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	vhostQueryRepo := NewVirtualHostQueryRepo(persistentDbSvc)

	t.Run("Read", func(t *testing.T) {
		withMappings := true
		_, err := vhostQueryRepo.Read(dto.ReadVirtualHostsRequest{
			Pagination:   dto.PaginationUnpaginated,
			WithMappings: &withMappings,
		})
		if err != nil {
			t.Errorf("ExpectingNoErrorButGot: %v", err)
		}
	})
}
