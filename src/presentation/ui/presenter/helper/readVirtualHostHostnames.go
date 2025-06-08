package uiPresenterHelper

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
)

func ReadVirtualHostHostnames(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) ([]string, error) {
	vhostHostnames := []string{}
	virtualHostLiaison := liaison.NewVirtualHostLiaison(persistentDbSvc, trailDbSvc)

	vhostResponseLiaisonOutput := virtualHostLiaison.Read(map[string]interface{}{
		"itemsPerPage": 1000,
		"withMappings": false,
	})
	if vhostResponseLiaisonOutput.Status != liaison.Success {
		return vhostHostnames, errors.New("ReadVirtualHostLiaisonBadResponse")
	}

	vhostReadResponse, assertOk := vhostResponseLiaisonOutput.Body.(dto.ReadVirtualHostsResponse)
	if !assertOk {
		return vhostHostnames, errors.New("AssertReadVirtualHostsResponseFailed")
	}

	for _, vhostEntity := range vhostReadResponse.VirtualHosts {
		if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
			continue
		}

		vhostHostnames = append(vhostHostnames, vhostEntity.Hostname.String())
	}

	return vhostHostnames, nil
}
